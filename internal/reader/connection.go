package reader

type ConnectionPool interface {
	Init([]string) error
}
type Connection interface{}

//GetResources(aliases, type, ?filter)
//GetRelatedResources(aliases, base_type, target_type, ?filter)
//IsSupported(type, ?attrs, ?filters)
//IsRelated(base_type, target_type)

// GetResources(RES, ATTRS, FILTERS) -> GetRelatedResources(response_item, RES2, ATTRS2, FILTERS2)

// GetResources(type: Type, attr: list, filter: lambda?, parent?: (type: str, id:str, relation:str))

// - Attrs() - return attribute names and its aliases can be specified in query
// - Addable() - return resources and its aliases that can be added
// - Pull(RESOURCE, ATTRIBUTES, FILTER, [PARENT]) - pulls the given resource with some attributes and some filter and an optional parent resource for better context
// - Resources() - lists supported resources and its aliases
// - Name() - returns plugin name, such as "aws"
