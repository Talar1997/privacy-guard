package blocker

type Blocker interface {
	SetRule(tvAddress string)
	RemoveRule(tvAddress string)
}
