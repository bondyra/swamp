package topology

type linkSchema struct {
	Links []linkJson `json:"links,omitempty" validate:"dive"`
}

type linkJson struct {
	From    NamespacedType `json:"from" validate:"required"`
	To      NamespacedType `json:"to" validate:"required"`
	Mapping mappingJson    `json:"mapping" validate:"required"` // todo many mappings
}

type mappingJson struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}
