package java

type Method interface {
	MethodPrototype() string
	MethodBody(returning Interface) string
}
