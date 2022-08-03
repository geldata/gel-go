quality: lint test bench

lint:
	golangci-lint run --sort-results

test:
	go test -v -race -bench=$$^ -timeout=10m ./...

bench:
	go test -run=^$$ -bench=. -benchmem -timeout=10m ./...

format:
	gofmt -s -w .

errors:
	type edb || (\
		echo "the edb command must be in your path " && \
		echo "see https://www.edgedb.com/docs/internals/dev/#building-locally" && \
		exit 1 \
		)
	edb gen-errors-json --client | \
		go run internal/cmd/generr/*.go > \
		generatederrors.go
	make format
