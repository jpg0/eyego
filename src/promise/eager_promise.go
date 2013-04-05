package promise

type EagerPromise struct {
	completed chan interface {}
}

func NewEagerPromise(f func () interface {}) EagerPromise {
	rv := &EagerPromise{make(chan interface {}, 1)}
	go func(){
		rv.completed <- f()
	}()
	return *rv
}


func (p EagerPromise) Get() interface {} {
	return <- p.completed
}

func (p EagerPromise) Then(next func (interface{}) interface {}) Promise {
	rv := &EagerPromise{make(chan interface {}, 1)}
	go func(){
		rv.completed <- next(<- p.completed)
	}()
	return *rv
}
