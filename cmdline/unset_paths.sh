export PATH="$(echo $PATH|sed -r 's%'$(dirname $(readlink -f ${BASH_SOURCE[0]}))'/bin:%%g')"
export GOPATH=""