package definition

import (
	"fmt"
	"os"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/bondyra/swamp/internal/reader"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

var validate *validator.Validate

type Factory interface {
	FromFile(string) (*Definition, error)
}

type DefaultFactory struct{}

func (df DefaultFactory) FromFile(jsonPath string) (*Definition, error) {
	var err error
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	return common.Unmarshal[Definition](data)
}

type DefinitionInterface interface {
	Validate() error
}

type Definition struct {
	TypeDefinitions []TypeDefinition `json:"types" validate:"required,unique=Type,unique=Alias,dive"`
}

type TypeDefinition struct {
	Type            string             `json:"type" validate:"required"`
	IdentifierField string             `json:"identifierField" validate:"required"`
	Alias           string             `json:"alias" validate:"required"`
	Parents         []ParentDefinition `json:"parents,omitempty" validate:"unique=Type,dive"`
	Attrs           []Attr             `json:"attrs,omitempty" validate:"dive"`
}

type ParentDefinition struct {
	Type     string `json:"type" validate:"required"`
	LinkType string `json:"linkType" validate:"required,oneof=inline resourceModel"`
	Links    []Link `json:"links" validate:"required"`
}

type Link struct {
	ParentField string `json:"parentField" validate:"required"`
	Field       string `json:"field" validate:"required"`
}

type Attr struct {
	Field   string `json:"field" validate:"required"`
	IsExtra bool   `json:"isExtra,omitempty"`
}

func (d *Definition) Validate() error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(d)
	if err != nil {
		return err
	}
	// more complicated validations that probably cannot be handled by validation package
	if err2 := d.areAllParentTypesDefined(); err2 != nil {
		return err2
	}
	return nil
}

func (d *Definition) areAllParentTypesDefined() error {
	var parentTypes = make([]string, 0)
	for _, typeDefinition := range d.TypeDefinitions {
		for _, parent := range typeDefinition.Parents {
			parentTypes = append(parentTypes, parent.Type)
		}
	}
	undefinedParentTypes := common.Difference(parentTypes, d.AllDefinedTypes())
	if len(undefinedParentTypes) > 0 {
		return fmt.Errorf("invalid definition, following parents are not defined: %v", undefinedParentTypes)
	}
	return nil
}

func (d *Definition) AllDefinedTypes() []string {
	allDefinedTypes := make([]string, 0)
	for _, typeDefinition := range d.TypeDefinitions {
		allDefinedTypes = append(allDefinedTypes, typeDefinition.Type)
	}
	return allDefinedTypes
}

func (d *Definition) GetAtributesForType(itemType string, all bool) []string {
	result := make([]string, 0)
	for _, td := range d.TypeDefinitions {
		if td.Type == itemType {
			for _, attr := range td.Attrs {
				if all || !attr.IsExtra {
					result = append(result, attr.Field)
				}
			}
		}
	}
	return result
}

func (d *Definition) SupportsType(itemType string) bool {
	return slices.Contains(d.AllDefinedTypes(), itemType)
}

func (d *Definition) SupportsAttrs(itemType string, attrs []string) bool {
	definedAttrs := d.GetAtributesForType(itemType, true)
	return len(common.Difference(attrs, definedAttrs)) == 0
}

func (d *Definition) SupportsFilter(itemType string, filter reader.Filter) bool {
	definedAttrs := d.GetAtributesForType(itemType, true)
	return slices.Contains(definedAttrs, filter.Attr)
}
