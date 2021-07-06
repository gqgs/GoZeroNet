package bigfile

import (
	"reflect"
	"testing"
)

func TestUnpackPieceField(t *testing.T) {
	tests := []struct {
		name   string
		packed PieceField
		want   string
	}{
		{
			"3, 6, 1",
			PieceField{3, 6, 1},
			"1110000001",
		},
		{
			"0, 9, 1",
			PieceField{0, 9, 1},
			"0000000001",
		},
		{
			"10",
			PieceField{10},
			"1111111111",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnpackPieceField(tt.packed); got != tt.want {
				t.Errorf("UnpackPieceField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackPieceField(t *testing.T) {
	tests := []struct {
		name     string
		unpacked string
		want     PieceField
	}{
		{
			"3, 6, 1",
			"1110000001",
			PieceField{3, 6, 1},
		},
		{
			"0, 9, 1",
			"0000000001",
			PieceField{0, 9, 1},
		},
		{
			"10",
			"1111111111",
			PieceField{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PackPieceField(tt.unpacked); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PackPieceField() = %v, want %v", got, tt.want)
			}
		})
	}
}
