export PATH="$HOME/.plenv/bin:$PATH"
eval "$(plenv init -)"

export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin

export PATH=$PATH:/usr/local/go/bin
export PATH="$PATH:$GOPATH/bin"

export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go

export GO15VENDOREXPERIMENT=1


export PYENV_ROOT="$HOME/.pyenv"
export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"
