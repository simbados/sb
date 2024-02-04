package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sb/internal/types"
	"sb/internal/util"
	"strings"
)

func BuildSandboxProfile(context *types.Context) string {
	profile := "; Generated by cli tool sb\n"
	profile += readBaseProfile(&context.Paths) + "\n"
	profile += fmt.Sprintf(`
; Allow running the target and node binary
(allow process-fork) ; needed for process-exec, has no filters
(allow process-exec file-read* (literal "%s"))
(allow file-read* (literal "%s"))`, context.Paths.BinaryPath, context.Paths.BinPath)
	profile += fmt.Sprintf(`
(allow file-read*
%v
)`, buildLiteralFolders(context.Paths.WorkingDir))
	if context.Config.SbConfig != nil {
		profile = write(context.Config.SbConfig, profile)
		profile = read(context.Config.SbConfig, profile)
		profile = readWrite(context.Config.SbConfig, profile)
		profile = process(context.Config.SbConfig, profile)
		profile = netIn(context.Config.SbConfig, profile)
		profile = netOut(context.Config.SbConfig, profile)
	}
	util.LogDebug(profile)
	return minifyProfile(profile)
}

func minifyProfile(profile string) string {
	noComments := regexp.MustCompile(";.*")
	noNewLines := regexp.MustCompile("\r?\n|\r")
	profile = noComments.ReplaceAllString(profile, "")
	profile = noNewLines.ReplaceAllString(profile, " ")
	return strings.Trim(profile, " ")
}

func process(config *types.SbConfig, profile string) string {
	if len(config.Process) > 0 {
		allPaths := ""
		for _, write := range config.Process {
			allPaths += fmt.Sprintf("	%v\n", write)
		}
		allPaths = strings.TrimRight(allPaths, "\n")
		profile += fmt.Sprintf(`
; allow-run, enable specifc permissions often used by external processes

; Allow IPC
(allow ipc*)
(allow iokit*) ; required by chrome
(allow mach*) ; required by chrome

; Allow app sending signals
(allow signal)

; Allow reading preferences
(allow user-preference-read)

; Allow entrypoints
(allow process-exec process-exec-interpreter file-read*
  (literal "/private/var/select/sh")
  (literal "/usr/bin/env")
  (literal "/bin/sh")
  (literal "/bin/bash")
  (literal "/usr/sbin")
)

(allow file-read*
  (regex #"^/Users/[^.]+/Library/Preferences/(.*).plist")
  (regex #"^/Library/Preferences/(.*).plist")
  (literal "/Library")
  (subpath "/dev")
  (subpath "/private/var") ; critical for certs, /private/var/select/sh and the like
  (subpath "/private/etc") ; openssl.cnf and the like
)

(allow file-write*
  (subpath "/dev")
)
; file-write: enabled
(allow process-exec file-read*
%v
)`, allPaths)
	}
	return profile
}

func write(config *types.SbConfig, profile string) string {
	if config.Write != nil && len(config.Write) > 0 {
		allPaths := ""
		for _, write := range config.Write {
			allPaths += fmt.Sprintf("	%v\n", write)
		}
		allPaths = strings.TrimRight(allPaths, "\n")
		profile += fmt.Sprintf(`
; file-write: enabled
(allow file-write*
%v
)`, allPaths)
	}
	return profile
}

func readWrite(config *types.SbConfig, profile string) string {
	if config.ReadWrite != nil && len(config.ReadWrite) > 0 {
		allPaths := ""
		for _, write := range config.ReadWrite {
			allPaths += fmt.Sprintf("	%v\n", write)
		}
		allPaths = strings.TrimRight(allPaths, "\n")
		profile += fmt.Sprintf(`
; file-read: enabled
(allow file-read*
%v
)
; file-write: enabled
(allow file-write*
%v
)
`, allPaths, allPaths)
	}
	return profile
}

func read(config *types.SbConfig, profile string) string {
	if config.Read != nil && len(config.Read) > 0 {
		allPaths := ""
		for _, write := range config.Read {
			allPaths += fmt.Sprintf("	%v\n", write)
		}
		allPaths = strings.TrimRight(allPaths, "\n")
		profile += fmt.Sprintf(`
; file-read: enabled
(allow file-read*
%v
)`, allPaths)
	}
	return profile
}

func netOut(config *types.SbConfig, profile string) string {
	if config.NetworkOutbound {
		profile += `
; allow-net-outbound: enabled
(allow network-inbound
  (local tcp)
  (local udp)
)
(allow network-outbound)
(allow network-bind)
(allow system-socket)`
		profile = addFallBackForCert(config, profile)
	}
	return profile
}

func addFallBackForCert(config *types.SbConfig, profile string) string {
	// If we do not allow other process we need some file read for network access
	if len(config.Process) == 0 {
		profile += `
(allow file-read*
  (regex #"^/Users/[^.]+/Library/Preferences/(.*).plist")
  (regex #"^/Library/Preferences/(.*).plist")
  (literal "/Library")
  (subpath "/dev")
  (subpath "/private/var") ; critical for certs, /private/var/select/sh and the like
  (subpath "/private/etc") ; openssl.cnf and the like
)

(allow file-write*
  (subpath "/dev")
)
`
	}
	return profile
}

func netIn(config *types.SbConfig, profile string) string {
	if config.NetworkInbound {
		profile += `
; allow-net-inbound: enabled
(allow network-bind network-inbound
  (local tcp)
  (local udp))
(allow system-socket)
(allow file-read* (literal "/private/var/run/resolv.conf"))
`
		profile = addFallBackForCert(config, profile)
	}
	return profile
}

func buildLiteralFolders(path string) string {
	allPaths := ""
	currentPath := ""
	paths := strings.Split(path, "/")
	for _, singlePath := range paths[1:] {
		currentPath += "/" + singlePath
		allPaths += "	(literal \"" + currentPath + "\")\n"
	}
	return strings.TrimRight(allPaths, "\n")
}

func readBaseProfile(path *types.Paths) string {
	content, err := os.ReadFile(filepath.Join(path.HomePath, ".sb-config", "default.sb"))
	handleFileReadError(err)
	return string(content)
}

func handleFileReadError(err error) {
	if err != nil {
		util.LogErr("Something went wrong while reading the base profile", err)
	}
}
