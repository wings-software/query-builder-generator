package java

type Void struct {
}

func (_ Void) InterfaceName() string {
	return "void"
}

func (_ Void) ReturnFromThis() string {
	return ""
}