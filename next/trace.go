package next

type Trace interface {
	Println(level int, args ...interface{})
	Println0(...interface{})
	Println1(...interface{})
	Println2(...interface{})
	Extend(string) Trace
}
