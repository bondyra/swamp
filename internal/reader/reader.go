package reader

import "github.com/bondyra/swamp/internal/common"

type Properties map[string]string

type ItemData struct {
	Identifier string
	Properties *Properties
}

type Item struct {
	Profile string
	Data    *ItemData
}

type Condition struct {
	Attr  string
	Op    common.Operator
	Value string
}

type Reader interface {
	GetNamespace() string
	GetItemSchemaPath() string
	GetLinkSchemaPath() string

	GetSupportedProfiles() []string

	GetItems(itemType string, profiles []string, attrs []string, conditions []Condition) ([]*Item, error)
}
