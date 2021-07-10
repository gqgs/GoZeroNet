package bigfile

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
	"github.com/stretchr/testify/require"
)

func TestMerkleRoot(t *testing.T) {
	sum256HexDigest := func(b []byte) string {
		digest := sha256.Sum256(b)
		return hex.EncodeToString(digest[:])
	}

	tests := []struct {
		name   string
		hashes []string
		hasher hashFunc
		want   string
	}{
		{
			"empty input",
			[]string{""},
			sum256HexDigest,
			"",
		},
		{
			"single hash",
			[]string{"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
			sum256HexDigest,
			"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		},
		{
			"multiple hashes",
			[]string{
				"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
				"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d",
				"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6",
				"18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4",
				"3f79bb7b435b05321651daefd374cdc681dc06faa65e374e38337b88ca046dea",
			},
			sum256HexDigest,
			"d71f8983ad4ee170f8129f1ebcdd7440be7798d8e1c80420bf11f1eced610dba",
		},
		{
			"sha512 hash",
			[]string{
				"9925a54fe7fe03488e4bbdeddef906f9353b763f7a1f483653360901b6c7e5bb",
				"a148621b234397a6b347f1e84c3b6094e13094e356942e87a1dc7cc647595216",
			},
			crypto.Sha512_256,
			"4a324083f7dbc6bce0ce91b2f4a9900fb02fff2bb271366ec832c1f2672af1bb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, merkleRoot(tt.hashes, tt.hasher))
		})
	}
}
