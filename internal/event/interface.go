package event

type Bus interface {
	Put(interface{}) error
}
