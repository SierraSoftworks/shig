package core

type Output interface {
	Printf(format string, a ...interface{})
	Println(a ...interface{})
}
