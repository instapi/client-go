.ONESHELL:

lint:
	@command -v golangci-lint > /dev/null || \
		GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

	golangci-lint run --exclude-use-default=false \
		--enable=govet \
		--enable=unused \
		--enable=structcheck \
		--enable=varcheck \
		--enable=ineffassign \
		--enable=deadcode \
		--enable=typecheck \
		--enable=interfacer \
		--enable=goconst \
		--enable=scopelint \
		--enable=golint \
		--enable=staticcheck \
		--enable=gosimple \
		--enable=unconvert \
		--enable=goconst \
		--enable=goimports \
		--enable=maligned \
		--enable=misspell \
		--enable=unparam \
		--enable=prealloc \
		./...

test:
	@go test -cover -failfast -race ./...
	go test -failfast -run=^$ -bench=. -benchmem ./...

update:
	@go get -u ./...
	go mod tidy
	go mod vendor
