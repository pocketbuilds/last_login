// A setup like this is useful for manual testing of your plugin.

package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	plugin "github.com/pocketbuilds/last_login"
)

func main() {
	app := pocketbase.New()

	(&plugin.Plugin{
		// test config will go here
		FieldName: "last_login",
	}).Init(app)

	err := app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
