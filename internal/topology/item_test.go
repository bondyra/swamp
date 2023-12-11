package topology

import (
	"testing"
)

func TestValidateItemModel(t *testing.T) {
	tests := []struct {
		name    string
		is      []itemJson
		wantErr bool
	}{
		{
			name:    "test empty",
			is:      []itemJson{},
			wantErr: false,
		},
		{
			name: "test full schema",
			is: []itemJson{
				{Type: NamespacedType{"r1", "t1"}, Attrs: []attrJson{{Field: "a1_1"}, {Field: "a1_2"}}},
				{
					Type:  NamespacedType{"r1", "t2"},
					Attrs: []attrJson{{Field: "a2_1"}, {Field: "a2_2"}},
				},
			},
			wantErr: false,
		},
		{
			name: "test error when types are duplicated",
			is: []itemJson{
				{Type: NamespacedType{"r1", "t1"}},
				{Type: NamespacedType{"r1", "t1"}},
			},
			wantErr: true,
		},
		{
			name: "test error when any type is empty",
			is: []itemJson{
				{Type: NamespacedType{"r1", ""}},
				{Type: NamespacedType{"r2", "t2"}},
			},
			wantErr: true,
		},
		{
			name: "test error when any reader is empty",
			is: []itemJson{
				{Type: NamespacedType{"", "t1"}},
				{Type: NamespacedType{"r2", "t2"}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModel(itemSchema{Items: tt.is}); (err != nil) != tt.wantErr {
				t.Errorf("validateModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
