# Reimplementation of node-safe with capability to sandbox every script with custom profiles
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
Sorted by priority
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
- [x] Vigilant mode, ask at the end to proceed with command and config
- [x] Possibility to add overall config json file to apply to all commands (discuss if good idea?) - Will be solved with extend keyword
- [x] Extend other config file, e.g. commonNode.json could be applied to npm.json, npx.json
  - [x] Look for "__extend__" key while parsing. Limit to 2 parent configs
- [x] Validate json config files, if anything can not be parsed return error (e.g. no array provided)
  - [x] No array/bool provided
  - [x] Config keys duplicated (can lead to bug that only first config is applied, disallow double config keys) - handled with schema validation
- [x] Add json schema
- [x] Option for only applying root or local config ```-c=local``` ```-c=root```
  - [x] Option for selecting one specific sandbox profile ```-c="/Users/test/.sb-config/npm.json```
- [ ] Better project structure with composition and no global state
- [ ] Extend README with complete instructions
- [ ] More testing for critical parts of the tool
- [ ] Add support for removing config files with sb
- [ ] Add TLDR; for README

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
you should put it into $HOME/.sb-config/\<name-of-your-binary\>.json  

2. Local Config: Sb will look for a .sb-config/ directory in your current directory and subdirectories.  
So if your current directory is ~/stuff/develop. Sb will look for a .sb-directory including parent directories up to your home folder.  
So it will look in ```/Users/<user>/stuff/.sb-config``` and ```/Users/<users>/.sb-config```. Uses the same structure as the global config  


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

## Configuration file
Now the structure of one such json file from the global and local configuration is as followed (npm.json):
```
{
    "__extends__": "[local]/npm.json",
    "__root-config__": {
      "read": ["~/temp"],
      "write": ["~/writeThisDir"],
      "net-out": false,
      "net-in": false
    },
    "install" : {
      "read": ["[wd]/**"],
      "write": ["[wd]/package.json"],
      "process": ["[bin]/curl"]
      "net-out": true,
      "net-in": true
    }
}
```
Explanation:  
`__extends__`: Allows extension of another config file, allows path specific identifiers see [Path identifiers](#path-identifiers).  
Must be string and does not support multiple extensions. The max extension count is currently **2** which might be changed in the future.  
`__root-config__`: Will be applied for all commands of a binary. E.g. npm install will trigger building of a sandbox profile with `__root-config__` and `install`
`install`: Will be applied if binary uses command install, e.g. npm install  
All configuration have the following object (all optional):
1. `write`: Array of paths - which should be writeable
2. `read`: Array of paths - which should be readable
3. `process`: Array of paths - which contains binaries that should be able to be spawned in subprocesses  
If your script needs curl than you should allow it here.
4. `net-in`: Boolean - if incoming network traffic should be allowed (e.g. local development web server)
5. `net-out`: Boolean - if outgoing network traffic should be allowed 
## Path identifiers
Path identifiers are special directories that can be used for specifying path in `write`, `read`, `process` arrays and `__extend__` string
Available identifiers are:
1. `[wd]`: working directory
2. `[home]`: home directory
3. `~`: home directory
4. `[target]`: directory of the target binary (`sb ls -a` runs ls which is the target binary and which could be located at /usr/bin)
5. `[bin]`: binary directory of system where most binaries are located (could be wrong so if you make sure a binary is whitelisted add the path to the binary)
6. `[local]`: local config directory (if exists). E.g. `/Users/<user>/someProject/.sb-config`
7. `[root]`: root config directory (if exists). E.g. `/Users/<user>/.sb-config`

## Globs
Following globs are supported for `write`, `read` and `process`.  
**Attention**: `__extends__` does not support globs, only identifiers
1. `*`: Should be used only for specifying group of files.  
`[wd]/*.js` only javascript files in working directory
2. `**`: (at the end of path): Includes subdirectories and everything included.  
`[wd]/**` allows everything inside working directory and subdirectories
3. `**`: (middle of path): Includes subdirectories.  
   `[wd]/**/hello` Allows paths such as `/Users/home/what/is/this/hello` if `/Users/home/` is working directory
4. `..`: Not allowed at beginning, because relative path is ambiguous.  
`[wd]/..` specify one directory above working directory

**One important distinction is that you can allow to read `/Users/home/hello` which does not mean that
the process can read any files at `/Users/home/hello`, it just means it can see the directory.  
If you want to allow files inside `hello` than you need to specify it with `/Users/home/hello/**`.**

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
- --config (-c): Specify which configuration files to use can be local, root or path
E.g. ```sb -c=root npm``` (uses only root config)  
E.g. ```sb -c=local npm``` (uses only local config)
E.g. ```sb -c=/Users/name/hello/npm.json npm``` (uses only config from path)


## FAQS
1. But sandbox-exec is deprecated and usage is discouraged  
A: Yes that is true, but most browsers rely on sandbox-exec and v1 sandbox profiles still work.  
There is currently no other way to sandbox scripts or cli tools.
