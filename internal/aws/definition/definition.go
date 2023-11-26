package definition

import (
	"fmt"
	"os"
	"regexp"

	"github.com/bondyra/swamp/internal/aws/common"
)

type Factory interface {
	NewDefinition(string) (Definition, error)
}

type DefaultFactory struct {
	jsonPath string
}

func (df DefaultFactory) NewDefinition() (*Definition, error) {
	var err error
	data, err := os.ReadFile(df.jsonPath)
	if err != nil {
		return nil, err
	}
	return common.Unmarshal[Definition](data)
}

type DefinitionInterface interface {
	Validate() error
}

type Definition struct {
	TypeDefinitions []TypeDefinition `json:"types"`
	allDefinedTypes []string
}

type TypeDefinition struct {
	Type            string             `json:"type"`
	IdentifierField string             `json:"identifierField"`
	Alias           string             `json:"alias"`
	Parents         []ParentDefinition `json:"parents,omitempty"`
	Attrs           []Attr             `json:"attrs,omitempty"`
}

type ParentDefinition struct {
	Type     string `json:"type"`
	LinkType string `json:"linkType"`
	Links    []Link `json:"links"`
}

type Link struct {
	ParentField string `json:"parentField"`
	Field       string `json:"field"`
}

type Attr struct {
	Field string `json:"field"`
}

func (d *Definition) Validate() error {
	var err error
	d.allDefinedTypes = d.AllDefinedTypes()
	validators := []func() error{d.areTypesFormatted, d.areTypesNotDuplicated, d.areAllParentTypesDefined}
	for _, validator := range validators {
		if err = validator(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Definition) areTypesFormatted() error {
	regex := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
	for _, t := range d.allDefinedTypes {
		if !regex.Match([]byte(t)) {
			return fmt.Errorf("invalid type: %v", t)
		}
	}
	return nil
}

func (d *Definition) areTypesNotDuplicated() error {
	duplicatedElements := common.DuplicatedElements(d.allDefinedTypes)
	if len(duplicatedElements) > 0 {
		return fmt.Errorf("invalid definition, following types were defined more than once: %v", duplicatedElements)
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
	undefinedParentTypes := common.Difference(parentTypes, d.allDefinedTypes)
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
