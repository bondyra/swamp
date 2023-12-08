package reader

type Properties *map[string]string

type ItemData struct {
	Identifier string
	Properties Properties
}

type Item struct {
	Profile string
	Data    *ItemData
}

type Filter func(i *Item) bool
type Transform func(props Properties) Properties

type Reader interface {
	Name() string

	GetSupportedProfiles() []string

	GetItems(itemType string, profiles []string, ids []string, filters []Filter, transforms []Transform) ([]*Item, error)
}
