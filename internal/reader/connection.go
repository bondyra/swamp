package reader

type ConnectionPool interface {
	Init([]string) error
}
type Connection interface{}
