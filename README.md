# The Go driver for Gel

[![Build Status](https://github.com/geldata/gel-go/workflows/Tests/badge.svg?event=push&branch=master)](https://github.com/geldata/gel-go/actions)
[![Join GitHub discussions](https://img.shields.io/badge/join-github%20discussions-green)](https://github.com/geldata/gel/discussions)

## Installation

In your module directory, run the following command.

```bash
$ go get github.com/geldata/gel-go
```

## Basic Usage

Follow the [Gel tutorial](https://www.geldata.com/docs/guides/quickstart)
to get Gel installed and minimally configured.

```go
package main

import (
	"context"
	"fmt"
	"log"

	gel "github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
)

func main() {
	ctx := context.Background()
	client, err := gel.CreateClient(ctx, gelcfg.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var result string
	err = client.QuerySingle(ctx, "SELECT 'hello Gel!'", &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
```

## Development

A local installation of Gel is required to run tests.
Download Gel from [here](https://www.geldata.com/download)
or [build it manually](https://www.geldata.com/docs/reference/dev).

To run the test suite run `make test`.
To run lints `make lint`.

## License

gel-go is developed and distributed under the Apache 2.0 license.
