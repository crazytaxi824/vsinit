package gocriticcheck

type Ints interface {
	~int | ~int32 | ~int64
}

type Byt interface {
	~byte
}

var length = 100
