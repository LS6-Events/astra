package otherpkg1

type Foo struct {
	Bar string
}

func GetFoo() (Foo, error) {
	return Foo{
		Bar: "bar",
	}, nil
}
