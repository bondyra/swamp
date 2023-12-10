package topology

import (
	"testing"

	"github.com/bondyra/swamp/internal/aws"
	"github.com/bondyra/swamp/internal/reader"
)

func TestReaderTopologyLoader(t *testing.T) {
	_, err := ReaderTopologyLoader([]reader.Reader{aws.AwsReader{}})()
	if err != nil {
		t.Errorf("reader topology error = %v", err)
		return
	}
}

func TestNewTopology(t *testing.T) {
	tests := []struct {
		name    string
		is      itemSchema
		ls      linkSchema
		wantErr bool
	}{
		{
			name: "test full topology",
			is: itemSchema{
				Items: []itemJson{
					{Type: NamespacedType{"r1", "t1"}},
					{Type: NamespacedType{"r1", "t2"}},
					{Type: NamespacedType{"r2", "t3"}},
				},
			},
			ls: linkSchema{
				Links: []linkJson{
					{From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r1", "t2"}},
					{From: NamespacedType{"r1", "t2"}, To: NamespacedType{"r2", "t3"}},
				},
			},
			wantErr: false,
		},
		{
			name:    "test when items and links are empty",
			is:      itemSchema{},
			ls:      linkSchema{},
			wantErr: false,
		},
		{
			name: "test error when items are invalid",
			is: itemSchema{
				Items: []itemJson{
					{Type: NamespacedType{"r1", "t1"}},
					{Type: NamespacedType{"r1", "t1"}},
				},
			},
			ls:      linkSchema{},
			wantErr: true,
		},
		{
			name: "test error when links are invalid",
			is: itemSchema{
				Items: []itemJson{
					{Type: NamespacedType{"r1", "t1"}},
					{Type: NamespacedType{"r1", "t2"}},
				},
			},
			ls: linkSchema{
				Links: []linkJson{
					{From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r1", "t1"}},
				},
			},
			wantErr: true,
		},
		{
			name: "test error when items and links are invalid",
			is: itemSchema{
				Items: []itemJson{
					{Type: NamespacedType{"r1", "t1"}},
					{Type: NamespacedType{"r1", "t1"}},
				},
			},
			ls: linkSchema{
				Links: []linkJson{
					{From: NamespacedType{"r1", "t1"}, To: NamespacedType{"r1", "t1"}},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newDefaultTopology([]*itemSchema{&tt.is}, []*linkSchema{&tt.ls})
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTopology() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
