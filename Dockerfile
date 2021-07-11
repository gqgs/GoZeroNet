FROM golang:1.16 AS builder
COPY . /app
WORKDIR /app

RUN go build -ldflags="-extldflags=-static" -o go-zeronet

FROM scratch
COPY --from=builder /app/go-zeronet /usr/bin/go-zeronet
COPY --from=builder /app/zeronet.toml /app/zeronet.toml

CMD ["/usr/bin/go-zeronet", "server", "--ui_server_addr",  "0.0.0.0:43111", "--file_server_addr", "0.0.0.0:26553"]

EXPOSE 43111 26553
