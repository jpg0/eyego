package promise

type Promise interface {
	Get() interface {};
	Then(func (interface{}) interface {}) Promise;
}
