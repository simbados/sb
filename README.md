# Reimplementation of node-safe with capability to sandbox every script with custom settings
- Reimplementation of node-safe (https://github.com/berstend/node-safe) in go
# Table of Content
1. [Status](#status)  
2. [TODOs](#todos)  
3. [Installation](#installation)  
4. [Difference to node-safe](#difference-to-node-safe)  
5. [How does it work](#how-does-it-work)
6. [Configuration Flags](#configuration-flags)
7. [FAQs](#faqs)

## Status

- Currently in development
- v0.1.0 coming soon

## TODOs
- [x] Glob support of ../
- [x] Add support for opening config files with sb and default editor
- [x] Implement init method
  - [x] Adds .sb-config default files
  - [x] Add default.sb to .sb-config root repo so that user can adjust default sandbox profile if needed
  - [x] Add install script to make life easier for installation and adding sb to path
  - [x] ~~Add sourcing of shell when adding to path (not possible because not main process in shell)~~
  - [x] Support multiple shells (fish, zsh, bash, etc)
- [x] ~~Flag for creating shim executable for binary (aliasing will be enough for now)~~
- [x] Fixing node binary call in npm.json (currently for nvm)
- [x] Sophisticated parsing of config json keys with glob expansion
- [x] Add support of adding new config files with sb
- [x] Helper command to show all configs file for a binary and their content and which will be applied
- [x] Print command that is run when passing args to sandbox-exec
- [x] Always apply __root-config__ of root config for binary and show in sb -s command
- [x] Merge local and root config
- [ ] Validate json config files, if anything can not be parsed return error (e.g. no array provided)
  - [x] No array/bool provided
  - [ ] Config keys duplicated (can lead to bug that only first config is applied, disallow double config keys)
- [ ] Add support for removing config files with sb
- [ ] Possibility to add overall config json file to apply to all commands (discuss if good idea?)
- [ ] Add TLDR; for README
- [x] Vigilant mode, ask at the end to proceed with command and config
- [ ] Option for only applying root or local config ```-c local -c root```
- [ ] More testing for critical parts of the tool

## Installation
**Building from Source**  

Run ```go build cmd/sb/sb.go && ./sb --init```  
This will build sb and run the [init](#2-configuration-flags) command which will setup everything for sb

**Prebuild binary**

Coming soon: With the Prebuild binary you can download sb executable directly from github and run the init function

## Difference to node-safe

- Generic approach for all binaries, node-safe is purposely build for node/npm/npx/yarn. You can also run other executables with node-safe,  
but sb lets you create own configurations for every binary.
- Sb has a root config which can be configured once and then you can forget about it
- Node-safe development seems to be stale
- Sb can accumulate default profiles for binaries and ship them out of the box (requires community effort)

## How does it work

There are 3 ways to configure sb  
None are mandatory, but if you do not provide any arguments the program will have nearly no permissions

1. Global Config: If you have a stuff that should always be applied to your config,
you should put it into $HOME/.sb-config/\<name-of-your-binary\>.json with the structure:  
\> .sbconfig  
&emsp; \> npm.json  
&emsp; \> ls.json  
Now the structure of one such json file is as followed (npm.json): 
```
{
	"__root-config__": {
		"read": ["~/temp"],
		"write": ["~/writeThisDir"],
		"net-out": false,
		"net-in": false
	},
	"install" : {
		"read": ["[wd]/**"],
		"write": ["[wd]/package.json"],
		"net-out": true,
		"net-in": true
	}
}
```
`__root-config__` will be applied for all commands of a binary, e.g. npm install will trigger building of a sandbox  
profile with `__root-config__` and `install`
2. Local Config: Sb will look for a .sb-config/ directory in your current directory and subdirectories.  
So if your current directory is ~/stuff/develop  
Sb will look for .sb-directory in subdirectories up till your home folder.  
So it will look in /Users/\<user\>/stuff/.sb-config and /Users/\<users\>/.sb-config  
Uses the same structure as the global config  
3. Cli options: You can pass the options also as cli flags  
E.g. 
- ```sb --read="[wd]/**" ls -a```
- ```sb --read="[wd]/**,~/what/yes" ls -a```

**The config files need arrays for read/write/read-write/process, the cli config expects a string
which is comma seperated!**

**Attention**: If you have multiple configs for a binary, sb will merge these configs, so it will append 
any strings from read/write/read-write/process and net out and net in will only be allowed if
all configs allow it. If one config does not allow it, then it will be forbidden.
E.g. root config does allow ```"net-out": true```, but your local config has ```"net-out": false```,
then net-out will be forbidden.

## Configuration flags

Run sb --help for all available flags  
- --debug (-d):  Show sandbox profile and debug information (will run the command)  
- --dry-run (-dr): Do not run sandbox just show debug information  
- --version (-v):  Show version  
- --init (-i): Will initialize sb. Add root configuration files and move sb to config binary location
- --edit (-e): Edit config files with your default editor (standard vim)   
Takes two additional arguments first must be either local/root and second name of binary you would like to edit
E.g. ```sb -e root npm``` (edits the npm.json file in the root directory)  
E.g. ```sb -e local npm``` (edits the npm.json file in the local directory)  
- --show (-s): Displays which config will be applied for binary, use as ```sb -s npm```


## FAQS
1. But sandbox-exec is deprecated and usage is discouraged  
A: Yes that is true, but most browsers rely on sandbox-exec and v1 sandbox profiles still work.  
There is currently no other way to sandbox scripts or cli tools.
