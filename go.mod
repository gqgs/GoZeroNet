module github.com/gqgs/go-zeronet

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/bytedance/sonic v0.0.0-20210626024310-4e02b5b2fefb
	github.com/ethereum/go-ethereum v1.10.4
	github.com/fasthttp/websocket v1.4.3
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/cors v1.2.0
	github.com/golang/mock v1.3.1
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/minio/simdjson-go v0.2.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	github.com/zeebo/bencode v1.0.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1
)

replace github.com/antonfisher/nested-logrus-formatter => github.com/gqgs/nested-logrus-formatter v1.3.2-0.20210613163930-09740a972bc4
