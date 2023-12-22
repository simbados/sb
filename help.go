package main

import (
	"fmt"
	"os"
)

func printHelp() {
	fmt.Printf(
		"Welcome to the sandbox cli tool - short sb \n" +
			"The cli options are the following: \n" +
			"--debug -d: 		Debug the different steps in the application\n" +
			"--print -p: 		Print the sandbox profile and infos about the sandbox profile (e.g. which config files were loaded) \n" +
			"--dry-run -dr		Do not execute the sandbox with the wanted binary will set debug and print to true so that you have all the information what sb would do\n" +
			"--create-exe -ce	Creates executable shim for binary, so that you can use the sandboxed binary within other programs which would natively use the normal binary\n" +
			"--help -h		Print this help section\n" +
			"--version -v		Show which version of sb is installed\n")
	os.Exit(0)
}

func showVersion() {
	// TODO print version
	fmt.Printf("Not implemented, try later versions\n")
	os.Exit(0)
}
