# osx-env-sync

Synchronize OS X environment variables for command line and GUI applications from a single source

## Introduction

On OS X, command line applications and GUI applications are treated differently. (Well this can be put in a more technically correct manner but this is what you experience from a user's point of view.) One fundamental difference is that although it's straightforward to feed command line applications with environment variables, it's not so for GUI applications. It's even harder to feed both types of applications from a single source of definitions. (Well, one particular *workaround* is launching GUI applications from command line.) Moreover, OS X's relevant means for setting up environment variables (or initializing programs in general) have been changing over time in consecutive releases which makes the situation worse. [Hundreds of topics at Stack Overflow](http://stackoverflow.com/search?q=environment-variables+osx) is the living proof of this *mess*.

**osx-env-sync** provides a simple and effective solution for synchronizing environment variables for both command line and GUI applications from a single source. Unlike many other solutions/workarounds provided at Stack Overflow and many blogs, **osx-env-sync** is simple and works with the latest version of OS X. (Tested on version 10.10.5.)

## Background

The original version (from [ersiner/osx-env-sync](https://github.com/ersiner/osx-env-sync) parsed your `~/.bash_profile`
for lines starting with `export` and used those to run `launchctl setenv` against each line discovered.

This version runs a `bash` or `zsh` shell in login + interactive mode, and gets that to run the `env` command.
It parses the output of this into a version of your environment, and runs `launchctl setenv` on each env var.

There are also two options:  one written in Go (in the `osx-env-sync/` sub-dir) and another written in Ruby.

The Go version is more-correct in a lot of ways, but requires Go to be installed on your machine to compile the
code.
The Ruby version should be able to use the version of Ruby that's installed in macOS as standard -- but there's
no reason a more recent version shouldn't work with minimal tweaking.

## Installation

You can run the `make` in the root dir, that will run against the included `Makefile`.
This will install the Go version, if Go is installed, otherwise
fall-back to the Ruby version.

It will assume the `$SHELL` env-var contains your preferred shell, but will then fall-back to using `zsh`.
The setting can be further overridden in a config file.

### Config file

The Go version can parse a Config file named `.osx-env-sync.toml` and read a few items from that:

- debug -- Whether to produce debug output
- noop -- Don't *actually* run the commands, therefore don't apply any changes
- shell -- Override the default Shell

## Updating

And run the script whenever you want to reload your environment variables:

`osx-env-sync-now`

Finally make sure that you relaunch your already running applications (including Terminal.app) to make them aware of the changes.

*The setup is persistent. It will survive restarts and relogins.*

## Details

On OS X, each Terminal.app session (window or tab) is accompanied with your configured login shell, `/bin/bash` by default.
During `bash` login shell startup, the following files are sourced in order:

- `/etc/profile`
- `/etc/bashrc`
- `~/.bash_profile`

During `zsh` login, it's:

- `/etc/zshenv`
- `~/.zshenv`
- `/etc/zprofile`
- `~/.zprofile`
- `/etc/zshrc`
- `~/.zshrc`
- `/etc/zlogin`
- `~/.zlogin`

NB: The more-complicated style from `zsh` was part of why I rewrote this to extract it from a shell that was running in login + interactive mode.

These files are sourced upon user login as well. But before these files, there is another type of script (well, there are
others too) executed on behalf of the user: Launch Agents.
**osx-env-sync** provides a launch agent which initializes a login shell and then uses `launchctl` command to set
environment variables for the whole user session.
A script helps the launch agent by parsing `~/.bash_profile` *just for reading the names of the environment variables*;
values of the variables are already effective in the launch agent execution environment as the shell script is run with
a *login shell* by the agent.

### More Details

While `/etc/profile` is being sourced (for `bash`) `/usr/libexec/path_helper` program is executed to set initial values
of `$PATH` and `$MANPATH` environment variables.
The program processes `/etc/paths` file as well as `/etc/paths.d/` and `/etc/manpaths.d/` directories for bootstrapping
the variables. You can also edit contents of these files and directories for system wide effect.
