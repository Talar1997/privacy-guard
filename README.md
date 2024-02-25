# Privacy-Guard

Application in Go that blocks outgoing communication from Sony TV. 

When TV is turned off (for example by clicking OFF button on tv remote) it goes to "Stand By" mode, which means device isn't actualy off. Display does not work but TV is able to communicate with internet and do other stuff. It also means it can talk with advertising services, send statistical data and so on.

Having privacy in mind, the application watch for status of TV, and if it's Stand By, it sets custom filtering rules on Adguard DNS to prevent sending any data while TV should not be working. When TV goes to "Active" state, it removes the rule to allow device work as expected.

Application is pretty small, it consumes ~6MB of RAM and works on linux alpine.

![Some of blocked request from TV in stand by mode, screen from Adguard Home](https://i.imgur.com/Yn1MKzm.png)
## Prerequisites:
- Running instance of Adguard Home (https://adguard.com/en/adguard-home/overview.html)
- Sony TV that expose REST API (https://pro-bravia.sony.net/develop/index.html)

## TODO
- Support for PiHole
- Integration with other TVs

## Build
```
go build -v -o ./build/privacy-guard  ./src/main.go
```
or
```
docker build  -t privacy-guard .
```

## Test
```
go vet

go test -coverprofile cover.out -v ./src/...

go tool cover -html=cover.out -o=cover.html
```

## Run
Make sure you provided required variables (You can use default.env as template)

```
./privacy-guard /path/to/.env
```

or 

```
export TV_ADDRESS=http://192.168.1.2 \
export ADGUARD_ADDRESS=http://192.168.1.3 \
export ADGUARD_USERNAME=user \
export ADGUARD_PASSWORD=passwd \
export INTERVAL=2 \
./privacy-guard
```

```
export TV_ADDRESS=xyz; export ADGUARD_ADDRESS=1; ./privacy-guard
```

or 

```
docker run privacy-guard
```
