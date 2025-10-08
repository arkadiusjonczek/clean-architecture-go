#!/bin/bash

FULLPATH="$(realpath "$0")"
FULLPATH="$(dirname "$FULLPATH")"

docker run -it -v $PWD:$PWD plantuml/plantuml -tsvg $FULLPATH/docs/*.plantuml
