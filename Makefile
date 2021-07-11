server: build
	./go-zeronet server

build:
	go build -o go-zeronet

run: server

docker:
	docker build -t go-zeronet .