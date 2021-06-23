package uiwebsocket

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/lib/websocket"
	"github.com/stretchr/testify/require"
)

func Test_channeljoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		payload  string
		expected func(w *uiWebsocket)
	}{
		{
			"channel object",
			`{"cmd":"channelJoin","params":{"channels":["siteChanged","serverChanged"]},"id":1000000}`,
			func(w *uiWebsocket) {
				for _, channel := range []string{"siteChanged", "serverChanged"} {
					_, joined := w.channels[channel]
					require.True(t, joined)
				}
			},
		},
		{
			"channel array",
			`{"cmd":"channelJoin","params":["siteChanged",
			"serverChanged"],"wrapper_nonce":"c53efa2fc5fbdfac74f25eb0afc7bca94fc9546c0000342065c31743edcfcdbd","id":1}`,
			func(w *uiWebsocket) {
				for _, channel := range []string{"siteChanged", "serverChanged"} {
					_, joined := w.channels[channel]
					require.True(t, joined)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := new(uiWebsocket)
			ws.ID = safe.Counter()
			ws.channels = make(map[string]struct{})
			mockConn := websocket.NewMockConn(ctrl)
			mockConn.EXPECT().WriteJSON(gomock.Any())
			ws.conn = mockConn
			err := ws.channelJoin([]byte(tt.payload), Message{})
			require.NoError(t, err)
			tt.expected(ws)
		})
	}
}
