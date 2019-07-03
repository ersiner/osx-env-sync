# Set-up the local environment:
# . $HOME/.bashrc
. $HOME/.zshenv
. $HOME/.zshrc

# Using Ruby to avoid having to escape spaces (etc) in name or val:
env |ruby -ne '$_ =~ /^(.+)=(.+)$/; name, val = $1, $2; system "launchctl", "setenv", name, val'
