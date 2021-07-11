server: build
	./go-zeronet server

build:
	go build -o go-zeronet

run:
	go run main.go server

docker:
	docker build -t go-zeronet .