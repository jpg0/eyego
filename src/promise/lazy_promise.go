package promise

type LazyPromise struct {
	f func () interface {}
}

func NewLazyPromise(f func () interface {}) *LazyPromise {
	return &LazyPromise{f}
}

func (p *LazyPromise) Get() {
	p.f()
}

func (p *LazyPromise) Then(next func (interface{}) interface {}) *LazyPromise {
	return &LazyPromise{func () interface {} {
		return next(p.f())
	}}
}
