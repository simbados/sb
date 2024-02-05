package log

import (
	"encoding/json"
	"fmt"
	"os"
	"sb/internal/types"
)

var Logger LoggerType = ImplLogger{}

type LoggerType interface {
	LogErr(args ...any)
	LogWarn(args ...any)
	LogDebug(args ...any)
	LogInfoLn(args ...any)
	LogInfoSl(args ...any)
	LogHighlight(args ...any)
	PrettyJson(toPretty any) string
}

type ImplLogger struct{}

func (il ImplLogger) LogErr(args ...any) {
	fmt.Print(ColorRed)
	fmt.Println(args...)
	fmt.Print(ColorReset)
	os.Exit(1)
}

func (il ImplLogger) LogHighlight(args ...any) {
	fmt.Print(ColorGreen)
	fmt.Println(args...)
	fmt.Print(ColorReset)
}

func (il ImplLogger) LogWarn(args ...any) {
	fmt.Print(ColorOrange)
	fmt.Println(args...)
	fmt.Print(ColorReset)
}

func (il ImplLogger) LogDebug(args ...any) {
	if types.CliOptions.DebugEnabled {
		fmt.Println(args...)
	}
}

func (il ImplLogger) LogInfoLn(args ...any) {
	fmt.Println(args...)
}

func (il ImplLogger) LogInfoSl(args ...any) {
	fmt.Print(args...)
}

func (il ImplLogger) PrettyJson(toPretty any) string {
	prettyJson, err := json.MarshalIndent(toPretty, "", "    ")
	if err != nil {
		LogErr("Some error while pretty printing json", err)
		return ""
	}
	return string(prettyJson)
}

const (
	ColorRed    = "\033[38;5;1;1m"
	ColorGreen  = "\033[38;5;35;1m"
	ColorOrange = "\033[38;5;202;1m"
	ColorReset  = "\033[0m"
)

func LogErr(args ...any) {
	Logger.LogErr(args...)
}

func LogDev(args ...any) {
	if types.Envs.DevModeEnabled {
		fmt.Println(args...)
	}
}

func LogHighlight(args ...any) {
	Logger.LogHighlight(args...)
}

func LogWarn(args ...any) {
	Logger.LogWarn(args...)
}

func LogDebug(args ...any) {
	Logger.LogDebug(args...)
}

func LogInfoLn(args ...any) {
	Logger.LogInfoLn(args...)
}

func LogInfoSl(args ...any) {
	Logger.LogInfoSl(args...)
}

func PrettyJson(toPretty any) string {
	return Logger.PrettyJson(toPretty)
}
