package store

type ConnectionStore interface {
	List() []string
	Add(id string) error
	Remove(id string) error
}
