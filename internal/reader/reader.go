package reader

type Filter struct{}

type ParentContext struct{}

type ItemData struct{}

type Reader interface {
	// reader name, for namespace query validation
	GetReaderName() string
	// list profile names that were loaded
	GetProfileNames() []string
	// list all item names that can be read by this reader
	GetItemNames() []string
	// list default attributes item with given name would have
	GetDefaultItemAttributes(string) []string
	// list attributes that item with given name can have
	GetAllItemAttributes(string) []string

	GetItems(string, []string, Filter, ParentContext) ([]ItemData, error)
}
