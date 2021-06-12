package site

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseTrackerResponse(t *testing.T) {
	annouce, err := os.Open("testdata/announce")
	if err != nil {
		t.Fatal(err)
	}
	defer annouce.Close()

	parsed, err := parseTrackerResponse(annouce)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, parsed.Incomplete, 654)
	require.Equal(t, parsed.Interval, 1800)
	require.Len(t, parsed.Peers, 60)

	expectedPeerList := []string{
		"85.66.195.66:1",
		"101.224.116.189:1",
		"92.12.184.68:1",
		"67.140.213.194:1",
		"101.32.219.55:1",
		"2.44.154.216:1",
		"81.177.250.44:1",
		"185.220.101.3:1",
		"34.241.40.205:1",
		"185.220.100.243:1",
	}

	peerList, err := parsePeers(parsed.Peers)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expectedPeerList, peerList))
}
