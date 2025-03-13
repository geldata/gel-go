package main

import (
	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
)

func main() {
	_, _ = gel.CreateClient(gelcfg.Options{})
}
