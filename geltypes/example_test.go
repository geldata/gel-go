package geltypes_test

import (
	"fmt"
	"log"

	"github.com/geldata/gel-go/geltypes"
)

func ExampleOptional() {
	type User struct {
		geltypes.Optional
		Name string `gel:"name"`
	}

	var result User
	query := `
		SELECT User { name }
		FILTER .name = "doesn't exist"
		LIMIT 1
	`
	err := client.QuerySingle(ctx, query, &result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Missing())

	err = client.QuerySingle(ctx, `SELECT User { name } LIMIT 1`, &result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Missing())

	// Output:
	// true
	// false
}
