package repocsmoney

type ICsmoney interface {
}

type csmoney struct {
}

func NewCsmoney() ICsmoney {
	return &csmoney{}
}
