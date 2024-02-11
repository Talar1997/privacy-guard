package tv

type Status int

const (
	Off Status = iota
	StandBy
	Active
)

type Tv interface {
	GetStatus() Status
	GetAddress() string
}
