CHANGES:=$(shell git status --porcelain)

quality: lint test bench

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 run --sort-results

test:
	# temporarily skip edgeql-go tests until a gel-go version is published
	go test -v -count=1 -race -bench=$$^ -timeout=20m $(shell go list ./... | grep -v edgeql-go)

bench:
	go test -run=^$$ -bench=. -benchmem -timeout=10m ./...

format:
	gofmt -s -w .

errors:
	type edb || (\
		echo "the edb command must be in your path " && \
		echo "see https://www.edgedb.com/docs/guides/contributing/code#building-locally" && \
		exit 1 \
		)
	edb gen-errors-json --client | \
		go run internal/cmd/generrdefinition/main.go > internal/client/errors_gen.go
	edb gen-errors-json --client | \
		go run internal/cmd/generrexport/main.go > errors_gen.go
	make format

gen:
	go generate ./...

gendocs:
	go run internal/cmd/gendocs/*.go

gendocs-lint:
	go run internal/cmd/gendocs/*.go --lint
