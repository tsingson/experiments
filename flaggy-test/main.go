package main

import (
	"fmt"

	"github.com/integrii/flaggy"
)

var (
	// make a variable for the version which will be set at build time
	version = "unknown"

	// keep subcommands as globals so you can easily check if they were used later on
	mySubcommand *flaggy.Subcommand
)

func init() {
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("Test Program")
	flaggy.SetDescription("A little example program")

	// you can disable various things by changing bools on the default parser
	// (or your own parser if you have created one)
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// you can set a help prepend or append on the default parser
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/integrii/flaggy"

	// create any subcommands and set their parameters
	mySubcommand = flaggy.NewSubcommand("mySubcommand")
	mySubcommand.Description = "My great subcommand!"

	// set the version and parse all inputs into variables
	flaggy.SetVersion(version)
	flaggy.Parse()
}

func main() {
	if mySubcommand.Used {
		fmt.Println("ok")
	}
}
