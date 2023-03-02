package gocriticcheck

// typeDefFirst
func (f f) typeDefFirst() {}

type f struct{}

// ok ✅
func (f ff[T]) typeDefFirstGeneric() {}

type ff[T Ints] struct {
	t T
}
