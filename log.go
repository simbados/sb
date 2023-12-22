package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func logErr(args ...any) {
	fmt.Println(args...)
	os.Exit(1)
}

func logWarn(args ...any) {
	if CliOptions.DebugEnabled {
		fmt.Println(args)
	}
}

func logInfo(args ...any) {
	if CliOptions.DebugEnabled {
		fmt.Println(args)
	}
}

func prettyJson(context Context) string {
	prettyJson, err := json.MarshalIndent(context, "", "    ")
	if err != nil {
		logErr("Some error while pretty printing json", err)
		return ""
	}
	return string(prettyJson)
}
