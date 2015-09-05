# osx-env-sync

Synchronize OS X environment variables for command line and GUI applications from a single source

## Introduction

On OS X, command line applications and GUI applications are treated differently. (Well this can be put in a more technically correct manner but this is what you experience from a user's point of view.) One fundamental difference is that although it's straightforward to feed command line applications with environment variables, it's not so for GUI applications. It's even harder to feed both types of applications from a single source of definitions. (Well, one particular *workaround* is launching GUI applications from command line.) Moreover, OS X's relevant means for setting up environment variables (or initializing programs in general) have been changing over time in consecutive releases which makes the situation worse. [Hundreds of topics at Stack Overflow](http://stackoverflow.com/search?q=environment-variables+osx) is the living proof of this *mess*.

**osx-env-sync** provides a simple and effective solution for synchronizing environment variables for both command line and GUI applications from a single source. Unlike many other solutions/workarounds provided at Stack Overflow and many blogs, **osx-env-sync** is simple and works with the latest version of OS X. (Tested on version 10.10.5.)

## Usage

Make sure your `~/.bash_profile` has the necessary export statements. You can have more than one export statement for a variable and you can use variable substitution as well, in the same way you define environment variables as usual. Here is an example:

```
export JAVA_HOME="$(/usr/libexec/java_home -v 1.8)"
export GOPATH="$HOME/go"
export PATH="$PATH:/usr/local/opt/go/libexec/bin:$GOPATH/bin"
export PATH="/usr/local/opt/coreutils/libexec/gnubin:$PATH"
export MANPATH="/usr/local/opt/coreutils/libexec/gnuman:$MANPATH"
```

Open your terminal and follow the steps below.

Download the launch agent:

`curl https://raw.githubusercontent.com/ersiner/osx-env-sync/master/osx-env-sync.plist -o ~/Library/LaunchAgents/osx-env-sync.plist`

Download the shell script:

`curl https://raw.githubusercontent.com/ersiner/osx-env-sync/master/osx-env-sync.sh -o ~/.osx-env-sync.sh`

Make sure the shell script is executable:

`chmod +x ~/.osx-env-sync.sh`

Load the launch agent for current session:

`launchctl load ~/Library/LaunchAgents/osx-env-sync.plist`

(Re)Launch a GUI application and verify that it can read the environment variables.

*The setup is persistent. It will survive restarts and relogins.*

After the initial setup (that you just did), if you want to reflect any changes in your `~/.bash_profile` to your whole environment again, rerunning the `launchctl load ...` command won't perform what you want; instead you'll get a warning like the following:

`<$HOME>/Library/LaunchAgents/osx-env-sync.plist: Operation already in progress`

In order to reload your environment variables without going through the logout/login process do the following:

`launchctl unload ~/Library/LaunchAgents/osx-env-sync.plist`

`launchctl load ~/Library/LaunchAgents/osx-env-sync.plist`

Finally make sure that you relaunch your already running applications (including Terminal.app) to make them aware of the changes.

## Details

On OS X, each Terminal.app session (window or tab) is accompanied with your configured login shell, `/bin/bash` by default. During login shell startup, the following files are sourced in order:

- `/etc/profile`
- `/etc/bashrc`
- `~/.bash_profile`

These files are sourced upon user login as well. But before these files, there is another type of script (well, there are others too) executed on behalf of the user: Launch Agents. **osx-env-sync** provides a launch agent which initializes a login shell *-achieved by passing `-l` parameter to `bash`-* so that `~/.bash_profile` is sourced in the first place and uses `launchctl` command to set environment variables for the whole user session. A shell script helps the launch agent by parsing `~/.bash_profile` *just for reading the names of the environment variables*; values of the variables are already effective in the launch agent execution environment as the shell script is run with a *login shell* by the agent.

### More Details

While `/etc/profile` is being sourced, `/usr/libexec/path_helper` program is executed to set initial values of `PATH` and `MANPATH` environment variables. The program processes `/etc/paths` file as well as `/etc/paths.d/` and `/etc/manpaths.d/` directories for bootstrapping the variables. You can also edit contents of these files and directories for system wide effect.
