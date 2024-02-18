package blocker

type Blocker interface {
	SetRule(rule string)
	RemoveRule(rule string)
}
