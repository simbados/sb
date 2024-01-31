package util

import (
	"encoding/json"
	"fmt"
	"os"
	"sb/internal/types"
)

const (
	ColorRed    = "\033[38;5;1;1m"
	ColorOrange = "\033[38;5;202;1m"
	ColorReset  = "\033[0m"
)

func LogErr(args ...any) {
	fmt.Print(ColorRed)
	fmt.Println(args...)
	fmt.Print(ColorReset)
	os.Exit(1)
}

func LogWarn(args ...any) {
	fmt.Print(ColorOrange)
	fmt.Println(args...)
	fmt.Print(ColorReset)
}

func LogDebug(args ...any) {
	if types.CliOptions.DebugEnabled {
		fmt.Println(args...)
	}
}

func LogInfoLn(args ...any) {
	fmt.Println(args...)
}

func LogInfoSl(args ...any) {
	fmt.Print(args...)
}

func logDev(args ...any) {
	fmt.Println(args...)
}

func PrettyJson(context *types.Context) string {
	prettyJson, err := json.MarshalIndent(context, "", "    ")
	if err != nil {
		LogErr("Some error while pretty printing json", err)
		return ""
	}
	return string(prettyJson)
}
