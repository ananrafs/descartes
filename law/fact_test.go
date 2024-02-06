package law

import (
	"reflect"
	"testing"
)

func TestMakeFact_Generate(t *testing.T) {
	tests := []struct {
		name   string
		src    map[string]interface{}
		slug   string
		target map[string]interface{}
		want   map[string]interface{}
	}{
		{
			name: "empty target",
			src: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
			slug: "random-1",
			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeFact(tt.src).Generate(tt.slug)
			if !reflect.DeepEqual(got.Facts.GetMap(), tt.want) {
				t.Errorf("Generate() got = %v, want %v", got.Facts.GetMap(), tt.want)
			}
		})
	}
}
