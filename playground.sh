#!/usr/bin/env bash
if [ $# -ne 1 ]; then
	echo "Usage: $0 <dirname>" >&2
	exit 1
fi
stat "$1" > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Path $1 does not exist" >&2
	exit 2
fi
export PLAYGROUND_ROOT=$(dirname $(readlink -f ${BASH_SOURCE[0]}))
export PLAYGROUND=$(basename "$1")
TEMPRC=$(mktemp)
cat > $TEMPRC <<EOS
if [ -f \$HOME/.bash_profile ]; then
	. \$HOME/.bash_profile
elif [ -f \$HOME/.profile ]; then
	. \$HOME/.profile
elif [ -f \$HOME/.bashrc ]; then
	. \$HOME.bashrc
fi
export PS1="[$PLAYGROUND] \$PS1"
export PLAYGROUND_ROOT="$PLAYGROUND_ROOT"
export PLAYGROUND="$PLAYGROUND"
export GOPATH="$PLAYGROUND_ROOT/$PLAYGROUND"
export PATH="$PLAYGROUND_ROOT/$PLAYGROUND/bin:$PATH"

echo "************************************"
echo "* Now in Go playground $PLAYGROUND"
echo "************************************"
EOS

bash --rcfile $TEMPRC