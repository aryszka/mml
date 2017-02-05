package mml

func Query(o *Val, q *Val) *Val {
	switch ot := o.sys.(type) {
	case *structure:
		return ot.sys[q.sys.(string)]
	case *list:
		return ot.sys[q.sys.(int)]
	default:
		panic("unexpected")
	}
}
