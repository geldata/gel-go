package gel_test

import (
	"fmt"
	"log"

	gel "github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
)

type User struct {
	ID   geltypes.UUID `gel:"id"`
	Name string        `gel:"name"`
}

func ExampleClient() {
	db, err := gel.CreateClient(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert a new user.
	var inserted struct{ id geltypes.UUID }
	err = db.QuerySingle(
		ctx,
		`INSERT User { name := <str>$0 }`,
		&inserted,
		"Bob",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Select user.
	var user User
	err = db.QuerySingle(
		ctx,
		`
	SELECT User { name }
	FILTER .id = <uuid>$id
	`,
		&user,
		map[string]any{"id": inserted.id},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user.Name)
	// Output: Bob
}
