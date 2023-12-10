package topology

import (
	"testing"
)

func TestValidateLinkModel(t *testing.T) {
	tests := []struct {
		name    string
		ls      []linkJson
		wantErr bool
	}{
		{
			name:    "test empty",
			ls:      []linkJson{},
			wantErr: false,
		},
		{
			name: "test full schema",
			ls: []linkJson{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r2", "t2"},
					Mapping: mappingJson{From: "a1", To: "a2"},
				},
				{
					From: NamespacedType{"r2", "t2"}, To: NamespacedType{"r1", "t1"},
					Mapping: mappingJson{From: "a2", To: "a1"},
				},
			},
			wantErr: false,
		},
		{
			name: "test error when any type is empty",
			ls: []linkJson{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r2", ""},
					Mapping: mappingJson{From: "a1", To: "a2"},
				},
			},
			wantErr: true,
		},
		{
			name: "test error when any reader is empty",
			ls: []linkJson{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"", "t2"},
					Mapping: mappingJson{From: "a1", To: "a2"},
				},
			},
			wantErr: true,
		},
		{
			name: "test error when any mapping attr is empty",
			ls: []linkJson{
				{
					From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r2", "t2"},
					Mapping: mappingJson{From: "", To: "a2"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModel(linkSchema{Links: tt.ls}); (err != nil) != tt.wantErr {
				t.Errorf("validateModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
