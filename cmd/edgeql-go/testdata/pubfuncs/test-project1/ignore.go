package main

import (
	"context"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
)

func main() {
	_, _ = gel.CreateClient(context.Background(), gelcfg.Options{})
}
