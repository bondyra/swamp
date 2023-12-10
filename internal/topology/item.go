package topology

type itemSchema struct {
	Items []itemJson `json:"items,omitempty" validate:"unique=Type,dive"`
}

type itemJson struct {
	Type  NamespacedType `json:"type" validate:"required"`
	Attrs []attrJson     `json:"attrs,omitempty" validate:"dive"`
}

type attrJson struct {
	Field   string `json:"field" validate:"required"`
	IsExtra bool   `json:"isExtra,omitempty"`
}
