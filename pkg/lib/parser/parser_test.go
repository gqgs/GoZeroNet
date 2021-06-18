package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFixJSONSpacing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"given valid JSON",
			`{"address":"1BLogC9LN4oPDcruNz3qo1ysa133E9AGg8","ignore":".*","inner_path":"data/users/content.json","modified":1562293663}`,
			`{"address": "1BLogC9LN4oPDcruNz3qo1ysa133E9AGg8", "ignore": ".*", "inner_path": "data/users/content.json", "modified": 1562293663}`,
		},
		{
			"given valid JSON with empty object",
			`{"address":"1BLogC9LN4oPDcruNz3qo1ysa133E9AGg8","files":{},"ignore":".*","inner_path":"data/users/content.json","modified":1562293663}`,
			`{"address": "1BLogC9LN4oPDcruNz3qo1ysa133E9AGg8", "files": {}, "ignore": ".*", "inner_path": "data/users/content.json", "modified": 1562293663}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := FixJSONSpacing(reader)
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}
