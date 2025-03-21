package geltypes_test

import (
	"context"
	"log"
	"testing"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/internal/testserver"
)

var (
	// Define identifiers that are used in examples, but keep them in their own
	// file. This way the definition is not included in the example. This keeps
	// examples concise.
	ctx    context.Context
	client *gel.Client
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	opts := testserver.Options()

	var err error
	client, err = gel.CreateClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Execute(ctx, `
		START MIGRATION TO {
			module default {
				type User {
					required name: str {
						default := 'default';
					};
				};
				type Product {};
			}
		};
		POPULATE MIGRATION;
		COMMIT MIGRATION;

		INSERT User { name := 'Stephanie' };
	`)
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}
