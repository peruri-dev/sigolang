build:
	CGO_ENABLED=0 go build -o ./sigolang ./main.go

test:
	go test -v ./...

run:
	go run main.go

mock:
	mockery --all

start: build
	./sigolang | tee sigolang.log
