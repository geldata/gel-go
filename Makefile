CHANGES:=$(shell git status --porcelain)

quality: lint test bench

lint:
	# GREP_FOR_THIS_STRING_WHEN_CHANGING_GOLANGCI_LINT_VERSION
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 config verify
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 run --sort-results

test:
	# The test server shuts it self down after 10 seconds of inactivity.  If
	# the server is still running when a subsequent test run is started the old
	# server will be reused.  This saves developer time at the expense of
	# repeatability.  You can set EDGEDB_SERVER_AUTO_SHUTDOWN_AFTER_SECONDS to
	# extend the timeout. Use make kill-test-server to stop the currently
	# running test server.
	go test -v -count=1 -race -bench=$$^ -timeout=20m ./...

kill-test-server:
	kill $(shell jq -r '.pid' /tmp/edgedb-go-test-server-info)

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
		go run internal/cmd/generrconst/main.go > gelerr/errors_gen.go
	edb gen-errors-json --client | \
		go run internal/cmd/generrtype/main.go > internal/gelerr/errors_gen.go
	make format

gen:
	go generate ./...
