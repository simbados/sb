package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ColorRed    = "\033[38;5;1;1m"
	ColorOrange = "\033[38;5;202;1m"
	ColorYellow = "\033[38;5;11;1m"
	ColorReset  = "\033[0m"
)

func logErr(args ...any) {
	fmt.Print(ColorRed)
	fmt.Println(args...)
	fmt.Print(ColorReset)
	os.Exit(1)
}

func logWarn(args ...any) {
	if CliOptions.DebugEnabled {
		fmt.Print(ColorOrange)
		fmt.Println(args)
		fmt.Print(ColorReset)
	}
}

func logInfo(args ...any) {
	if CliOptions.DebugEnabled {
		fmt.Println(args)
	}
}

func logDev(args ...any) {
	fmt.Println(args)
}

func prettyJson(context *Context) string {
	prettyJson, err := json.MarshalIndent(context, "", "    ")
	if err != nil {
		logErr("Some error while pretty printing json", err)
		return ""
	}
	return string(prettyJson)
}
