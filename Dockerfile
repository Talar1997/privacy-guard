# Build
FROM golang:1.21.6 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /privacy-guard  ./src/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -coverprofile cover.out -v ./src/...

# Deploy
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /privacy-guard /privacy-guard

USER nonroot:nonroot

ENTRYPOINT ["/privacy-guard"]