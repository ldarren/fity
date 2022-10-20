start:
	./fity -config conf.json -addr :9999

run:
	go run main.go -config conf.json

build:
	go build

test:
	go test -v ./...
