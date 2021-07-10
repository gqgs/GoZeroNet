server: build
	./go-zeronet server

build:
	go build -o go-zeronet

run: server