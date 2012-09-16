this=$(dirname $(readlink -f ${BASH_SOURCE[0]}))
export GOPATH="$this"
export PATH="$this/bin:$PATH"