GOPATH_BIN:=$(shell go env GOPATH)/bin

build:
	CGO_ENABLED=0 go build -o ./sigolang ./main.go

test:
	go test -v ./...

lint:
	go vet ./... && golangci-lint run && gocyclo -over=15 .

run:
	go run main.go

mock:
	mockery

start: build
	./sigolang | tee sigolang.log

prepare-ci:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	go install github.com/t-yuki/gocover-cobertura@latest
	go install github.com/firodj/nocov@latest
	go install github.com/jstemmer/go-junit-report/v2@latest

test-ci:
	go test -json -cover -coverprofile=test-coverage.txt.tmp -covermode count -v ./... 2>&1 | tee test-result.json | $(GOPATH_BIN)/gotestfmt
	$(GOPATH_BIN)/nocov -coverprofile test-coverage.txt.tmp -ignore-gen-files -ignore-files "^cmd/" > test-coverage.txt
	$(GOPATH_BIN)/gocover-cobertura < test-coverage.txt > test-coverage.xml
	cat test-result.json | $(GOPATH_BIN)/go-junit-report -parser gojson -set-exit-code > test-report.xml
	go tool cover -func=test-coverage.txt
	go tool cover -html=test-coverage.txt -o test-coverage.html
	sed -i "s=<source>.*</source>=<source>./</source>=g" test-coverage.xml
	sed -i "s;filename=\"sigolang/;filename=\";g" test-coverage.xml

lint-ci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_BIN) v1.64.6
	$(GOPATH_BIN)/golangci-lint run
