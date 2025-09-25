#!/bin/bash

docker run -it -v $PWD:$PWD plantuml/plantuml -tsvg $PWD/*.diagram
