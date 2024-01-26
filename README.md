# Reimplementation of node-safe

- Currently in development
- Reimplementation of node-safe (https://github.com/berstend/node-safe) in go

Release tbd

# 1 General Idea

There are 3 ways to configure sb  
None are mandatory, but if you not provide any arguments the program will not be sandboxed  

1. Global Config: If you have a stuff that should always be applied to your config,
you should put it into $HOME/.sb-config/<name-of-your-binary>.json
2. Local Config: Sb will look for a .sb-config/ directory in your current directory.
Uses the same structure as the global config  
3. Cli options: You can pass the options also as cli flags. More on this in section 3  

Attention: If you have both configs for a binary, sb will take the arguments with the highest priority  
Order of priority: cli arguments > local config > global config
E.g. You deny net-inbound on the global config, then it will always be denied, even though your local config can allow it

So for 1 and 2 you have one config directory at $HOME/.sb-config and <workdir>/.sb-config:

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
