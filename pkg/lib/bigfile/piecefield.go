package bigfile

import (
	"strings"
)

type PieceField []int

// Unpacks a piecefield according to the documentation bellow:
// https://zeronet.io/docs/help_zeronet/network_protocol/#bigfile-piecefield
func UnpackPieceField(piecefield PieceField) string {
	var result string
	for i, count := range piecefield {
		value := "0"
		if i%2 == 0 {
			value = "1"
		}
		result += strings.Repeat(value, count)
	}
	return result
}

// Packs a piecefield according to the documentation bellow:
// https://zeronet.io/docs/help_zeronet/network_protocol/#bigfile-piecefield
func PackPieceField(piecefield string) PieceField {
	var result PieceField
	value := '1'
	var count int
	for _, i := range piecefield {
		if i == value {
			count++
			continue
		}
		if value == '1' {
			value = '0'
		} else {
			value = '1'
		}

		result = append(result, count)
		count = 1
	}
	result = append(result, count)
	return result
}
