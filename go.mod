module github.com/gqgs/go-zeronet

go 1.16

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/fasthttp/websocket v1.4.3
	github.com/go-chi/chi v1.5.4
	github.com/mailru/easyjson v0.7.7
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.3.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	github.com/zeebo/bencode v1.0.0
)

replace github.com/antonfisher/nested-logrus-formatter => github.com/gqgs/nested-logrus-formatter v1.3.2-0.20210613163930-09740a972bc4
