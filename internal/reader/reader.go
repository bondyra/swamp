package reader

type Credentials struct {
}

type ReadContext struct {
	parentType string
	parentId   string
	parentData any
}

type Item interface{}

type Reader interface {
	QueryAllProfiles() ([]Credentials, error)
	QueryProfiles(profileNames []string) ([]Credentials, error)
}

type ReaderExt interface {
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

	// for specified resource type, list of attrs, userFilter string (to evaluate) and optional ReadContext,
	ReadItems(string, []string, string, ReadContext) []Item
}

// - Pull(RESOURCE, ATTRIBUTES, FILTER, [PARENT]) - pulls the given resource with some attributes and some filter and an optional parent resource for better context
