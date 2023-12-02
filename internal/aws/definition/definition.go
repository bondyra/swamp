package definition

import (
	"fmt"
	"os"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func FromFile(jsonPath string) (*Definition, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("NewDefinition: %w", err)
	}
	return common.Unmarshal[Definition](data)
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

func Validate(d Definition) error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(d)
	if err != nil {
		return err
	}
	return nil
}
