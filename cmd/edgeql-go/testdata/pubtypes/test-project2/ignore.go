package main

import (
	"context"

	_ "test/argnames"
	_ "test/object"
	_ "test/scalar"
	_ "test/tuple"
	_ "test/when_no_go_files_in_dir_dir_name_becomes_package_name"
	_ "test/when_no_go_files_in_dir_dir_name_becomes_package_name/subpkg"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
)

func main() {
	_, _ = gel.CreateClient(context.Background(), gelcfg.Options{})
}
