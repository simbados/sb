package util

import (
	"fmt"
	"os"
)

func PrintHelp() {
	fmt.Printf(
		`Welcome to the sandbox cli tool - short sb
The cli options are the following:
--debug -d: 		Debug the different steps in the application
--print -p: 		Print the sandbox profile and infos about the sandbox profile (e.g. which config files were loaded)
--dry-run -dr		Do not execute the sandbox with the wanted binary will set debug and print to true so that you have all the information what sb would do
--help -h		Print this help section
--version -v		Show which version of sb is installed
--show -s		Show location of all config files for this binary and the content of the file that would be applied
--vigilant -vi		Print Profile and ask before running command`)
	os.Exit(0)
}
