package topology

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bondyra/swamp/internal/aws/common"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type NamespacedType struct {
	Namespace string `validate:"required"`
	Type      string `validate:"required"`
}

func (nt NamespacedType) String() string {
	return fmt.Sprintf("%s.%s", nt.Namespace, nt.Type)
}

func (a NamespacedType) Compare(b NamespacedType) int {
	return strings.Compare(a.String(), b.String())
}

func (nt *NamespacedType) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s := strings.Split(v.(string), ".")
	nt.Namespace = s[0]
	nt.Type = s[1]
	return nil
}

func loadFromFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loadFromFile: %w", err)
	}
	result, err := common.Unmarshal[T](data)
	if err != nil {
		return nil, fmt.Errorf("DefaultSchemaReader: %w", err)
	}
	err = validateModel(*result)
	if err != nil {
		return nil, fmt.Errorf("DefaultSchemaReader: %w", err)
	}
	return result, nil
}

func validateModel[T any](t T) error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(t)
	if err != nil {
		return err
	}
	return nil
}
