package common

import (
	"reflect"
	"testing"
)

func TestMapManipulator_Copy(t *testing.T) {
	tests := []struct {
		name   string
		src    map[string]interface{}
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
			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
		},
		{
			name: "empty source",
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
		},
		{
			name: "diff source & target",
			src: map[string]interface{}{
				"delta": 412,
				"beta":  "1455",
			},
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},

			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
				"delta": 412,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ManipulateMap(tt.src).Copy(tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapManipulator_DeepCopy(t *testing.T) {
	tests := []struct {
		name   string
		src    map[string]interface{}
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
			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
		},
		{
			name: "empty source",
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
		},
		{
			name: "diff source & target",
			src: map[string]interface{}{
				"delta": 412,
				"beta":  "1455",
			},
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},

			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
				"delta": 412,
			},
		},
		{
			name: "diff object source & target",
			src: map[string]interface{}{
				"delta": 412,
				"beta":  "1455",
			},
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
				"omega": map[string]interface{}{
					"3444": 3444,
					"3412": 3412,
				},
			},

			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
				"delta": 412,
				"omega": map[string]interface{}{
					"3444": 3444,
					"3412": 3412,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ManipulateMap(tt.src).DeepCopy(tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepCopy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapManipulator_Merge(t *testing.T) {
	tests := []struct {
		name   string
		src    map[string]interface{}
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
			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
		},
		{
			name: "empty source",
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},
		},
		{
			name: "diff source & target",
			src: map[string]interface{}{
				"delta": 412,
				"beta":  "1455",
			},
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
			},

			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
				"delta": 412,
			},
		},
		{
			name: "diff object source & target",
			src: map[string]interface{}{
				"delta": 412,
				"beta":  "1455",
				"omega": map[string]interface{}{
					"3412": 3412,
				},
			},
			target: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
				"omega": map[string]interface{}{
					"3444": 3444,
				},
			},

			want: map[string]interface{}{
				"alpha": 1312,
				"beta":  "1312",
				"gamma": "kk",
				"delta": 412,
				"omega": map[string]interface{}{
					"3444": 3444,
					"3412": 3412,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ManipulateMap(tt.src).Merge(tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() got = %v, want %v", got, tt.want)
			}
		})
	}
}
