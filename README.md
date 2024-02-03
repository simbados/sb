# Reimplementation of node-safe

- Currently in development
- Reimplementation of node-safe (https://github.com/berstend/node-safe) in go

# Difference to node-safe

- Generic approach for all binaries, node-safe is purposely build for node/npm/npx/yarn. You can also run other executables with node-safe,  
but sb lets you create own configurations for every binary.
- Sb has a root config which can be configured once and then you can forget about it
- Node-safe development seems to be stale
- Sb can accumulate default profiles for binaries and ship them out of the box (requires community effort)

# TODOs
- [ ] Implement init method
  - [x] Adds .sb-config default files
  - [x] Add default.sb to .sb-config root repo so that user can adjust default sandbox profile if needed
  - [x] Add install script to make life easier for installation and adding sb to path
  - [ ] ~~Add sourcing of shell when adding to path (not possible because not main process in shell)~~
  - [ ] Support multiple shells (fish, zsh, bash, etc)
- [ ] Flag for creating shim executable for binary
- [ ] Fixing node binary call in npm.json (currently for nvm)
- [ ] More testing for critical parts of the tool
- [ ] Sophisticated parsing of config json keys with glob expansion
- [ ] Glob support of ../ 


# 1. General Idea

There are 3 ways to configure sb  
None are mandatory, but if you not provide any arguments the program will not be sandboxed  

1. Global Config: If you have a stuff that should always be applied to your config,
you should put it into $HOME/.sb-config/<name-of-your-binary>.json  
Inside .sb-config you will have json files which are named after the binary, e.g:  
npm.json  
ls.json  
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
So it will look in /Users/<user>/stuff/.sb-config and /Users/<users>/stuff/.sb-config
Uses the same structure as the global config  
3. Cli options: You can pass the options also as cli flags. More on this in section 3  

**Attention**: If you have both configs for a binary, sb will take the arguments with the highest priority  
Order of priority: cli arguments > local config > global config
E.g. You deny net-inbound on the global config, but your local config allows it, then it will be allowed  


# FAQS
1. But sandbox-exec is deprecated and usage is discouraged  
A: Yes that is true, but most browsers rely on sandbox-exec and v1 sandbox profiles still work.  
There is currently no other way to sandbox scripts or cli tools.