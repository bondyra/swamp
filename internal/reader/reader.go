package reader

type Filter struct {
	Attr  string
	Op    string
	Value string
}

type ParentContext struct{}

type ItemData struct {
	Identifier string
	Properties *map[string]string
}

type Item struct {
	Profile    string
	Properties ItemData
}

type Reader interface {
	Name() string
	KnownTypes() []string

	IsTypeSupported(itemType string) bool
	IsLinkSupported(itemType string, parentType string) bool
	AreAttrsSupported(itemType string, attrs []string) bool
	IsFilterSupported(itemType string, filter Filter) bool

	GetItems(itemType string, profiles []string, attrs []string, filter Filter, context ParentContext) ([]*Item, error)
}
