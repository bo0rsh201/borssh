# borssh
Simple ssh wrapper that transfers your dot files via ssh according to config

# configuration
```
home dir: $HOME/.borssh
config file: $HOME/.borssh/config.toml
```
# config.toml:
*config paths are relative to $HOME/.borssh/ directory*
```
BashProfile = [ "bash/aliases", "bash/prompt" ]
```
# installation
You can get binaries [here](https://github.com/bo0rsh201/borssh/releases/latest)

Or install it manually with
```
go install github.com/bo0rsh201/borssh
```
*you should have Go [installed](https://golang.org/doc/install)*
# commands
```
borssh compile
```
compiles current version of config to:
- $HOME/.borssh/bash_profile.compiled - concat of all dotfiles from config
- $HOME/.borssh/hash.compiled md5 of previous file
and includes it into local .bash_profile

```
borssh connect <host>
```
checks remote config version, performs sync and install if required and does ssh

# flags
```
-q Quite mode (suppress all output)
```