package definition

import (
	"encoding/json"
	"os"
)

// TODO: validation

type Definition struct {
	TypeDefinitions []TypeDefinition `json:"types"`
}

type TypeDefinition struct {
	Type            string             `json:"type"`
	IdentifierField string             `json:"identifierField"`
	Alias           string             `json:"alias"`
	Parents         []ParentDefinition `json:"parents"`
	Attrs           []Attr             `json:"attrs"`
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

type Reader interface {
	ReadDefinition(string) (*Definition, error)
}

type DefaultReader struct {
}

func (dr *DefaultReader) ReadDefinition(jsonPath string) (*Definition, error) {
	output := Definition{}

	var err error
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}
