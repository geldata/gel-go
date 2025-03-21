package gelcfg_test

import (
	"fmt"
	"log"
	"time"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
)

func ExampleOptions() {
	opts := gelcfg.Options{
		ConnectTimeout:     60 * time.Second,
		WaitUntilAvailable: 5 * time.Second,
		WarningHandler:     gelcfg.WarningsAsErrors,
	}

	client, err := gel.CreateClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	var message string
	err = client.QuerySingle(ctx, `SELECT "hello Gel"`, &message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
	// Output: hello Gel
}
