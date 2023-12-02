package reader

type Filter struct {
	Attr  string
	Op    int
	Value string
}

const (
	OpEquals int = iota
	OpNotEquals
	OpLike
	OpNotLike
)

type ParentContext struct{}

type ItemData struct {
	Identifier string
	Properties *map[string]string
}

type Item struct {
	Profile string
	Data    *ItemData
}

type Reader interface {
	Name() string

	GetSupportedProfiles() []string
	IsTypeSupported(itemType string) bool
	IsLinkSupported(itemType string, parentType string) bool
	AreAttrsSupported(itemType string, attrs []string) bool
	IsFilterSupported(itemType string, filter Filter) bool

	GetItems(itemType string, profiles []string, attrs []string, filters []Filter) ([]*Item, error)
}
