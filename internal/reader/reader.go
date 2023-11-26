package reader

type Filter struct {
	Attr  string
	Op    string
	Value string
}

type ParentContext struct{}

type ItemData struct {
	Properties *map[string]string
}

type Reader interface {
	// reader name, for namespace query validation
	GetReaderName() string
	// list profile names that were loaded
	GetProfileNames() []string
	// list all item names that can be read by this reader
	GetItemTypes() []string

	GetItems(string, []string, Filter, ParentContext) ([]*ItemData, error)
}
