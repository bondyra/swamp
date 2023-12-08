package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type SchemaLoader func() (*Schema, error)

func DefaultSchemaLoader() SchemaLoader {
	return func() (*Schema, error) {
		_, filename, _, _ := runtime.Caller(0)
		return fromFile(path.Dir(filename) + "/schema.json")
	}
}

func fromFile(jsonPath string) (*Schema, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("DefaultSchemaReader: %w", err)
	}
	result, err := common.Unmarshal[Schema](data)
	if err != nil {
		return nil, fmt.Errorf("DefaultSchemaReader: %w", err)
	}
	err = Validate(*result)
	if err != nil {
		return nil, fmt.Errorf("DefaultSchemaReader: %w", err)
	}
	return result, nil
}

type Schema struct {
	Items []ItemSchema `json:"items,omitempty" validate:"unique=Type,dive"`
	Links []LinkSchema `json:"links,omitempty" validate:"dive"`
}

type ItemSchema struct {
	Type  NamespacedType `json:"type" validate:"required"`
	Attrs []Attr         `json:"attrs,omitempty" validate:"dive"`
}

type NamespacedType struct {
	Reader string `validate:"required"`
	Type   string `validate:"required"`
}

func (nt *NamespacedType) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s := strings.Split(v.(string), ".")
	nt.Reader = s[0]
	nt.Type = s[1]
	return nil
}

type Attr struct {
	Field   string `json:"field" validate:"required"`
	IsExtra bool   `json:"isExtra,omitempty"`
}

type LinkSchema struct {
	From     NamespacedType `json:"from" validate:"required"`
	To       NamespacedType `json:"to" validate:"required"`
	Mappings []Mapping      `json:"mappings" validate:"required,dive"`
}

type Mapping struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}

func Validate(s Schema) error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}

func (s Schema) IsTypeSupported(reader, typ string) bool {
	nt := NamespacedType{Reader: reader, Type: typ}
	for _, i := range s.Items {
		if i.Type == nt {
			return true
		}
	}
	return false
}

func (s Schema) IsLinkSupported(fromReader, fromType, toReader, toType string) bool {
	fnt := NamespacedType{Reader: fromReader, Type: fromType}
	tnt := NamespacedType{Reader: toReader, Type: toType}
	for _, l := range s.Links {
		if l.From == fnt && l.To == tnt {
			return true
		}
	}
	return false
}

func (s Schema) AreAttrsSupported(reader, typ string, attrs []string) bool {
	nt := NamespacedType{Reader: reader, Type: typ}
	for _, i := range s.Items {
		if i.Type == nt {
			supported := common.Map(i.Attrs, func(a Attr) string { return a.Field })
			return len(common.Difference(attrs, supported)) == 0
		}
	}
	return false
}
