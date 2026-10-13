GOPATH_BIN:=$(shell go env GOPATH)/bin
APP_VERSION:=$(shell git rev-parse --short HEAD)

build:
	CGO_ENABLED=0 \
		go build -ldflags="-X 'main.AppVersion=$(APP_VERSION)'" -o ./sigolang

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
	go install gotest.tools/gotestsum@latest
	go install github.com/t-yuki/gocover-cobertura@latest
	go install github.com/firodj/nocov@latest
	go install github.com/jstemmer/go-junit-report/v2@latest

migrate-ci:
	ENV=test go run main.go db:init
	ENV=test go run main.go db:migrate

test-ci-db:
	-$(GOPATH_BIN)/gotestsum --jsonfile test-result.json -- -tags integration -coverprofile=test-coverage.txt.tmp ./...

test-ci-unit:
	-$(GOPATH_BIN)/gotestsum --jsonfile test-result.json -- -coverprofile=test-coverage.txt.tmp ./...
	
test-ci: test-ci-unit
	-$(GOPATH_BIN)/nocov -coverprofile test-coverage.txt.tmp -ignore-gen-files -ignore-files "^cmd/" > test-coverage.txt
	-$(GOPATH_BIN)/gocover-cobertura < test-coverage.txt > test-coverage.xml
	-$(GOPATH_BIN)/go-junit-report -parser gojson -set-exit-code < test-result.json > test-report.xml
	go tool cover -func=test-coverage.txt
	go tool cover -html=test-coverage.txt -o test-coverage.html
	sed -i "s=<source>.*</source>=<source>./</source>=g" test-coverage.xml
	sed -i "s;filename=\"sigolang/;filename=\";g" test-coverage.xml

lint-ci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_BIN) v1.64.6
	$(GOPATH_BIN)/golangci-lint run

docker-build: app-version
	docker build . -t sigolang:latest

docker-run:
	docker run --rm --env-file=.env --network=host -it sigolang:latest

docker-sh:
	docker run --entrypoint=sh --rm --env-file=.env --network=host -it sigolang:latest

app-version:
	echo $(APP_VERSION) > APP_VERSION
