#!/bin/bash

# function to run go mod tidy in each directory

function tidy {
    for d in `ls -d */`
    do
        # remove the / from $d
        d=${d%?}
        go work use $d
        echo "Tidying $d"
        cd $d
        # check of go.mod exists and it contains the correct module name: module github.com/Zate/go-templates/$d
        if [ -f go.mod ]
        then
            if grep -q "module github.com/Zate/go-templates/$d" go.mod
            then
                go mod tidy
            else
                echo "go.mod exists but does not contain the correct module name"
                rm go.mod
                go mod init github.com/Zate/go-templates/$d
                go mod tidy
            fi
        else
            echo "go.mod does not exist"
            go mod init github.com/Zate/go-templates/$d
            go mod tidy
        fi
        cd ..
        go work sync
    done
}

# function to generate a new go template directory from called $1

function new {
    cp -rp base $1
    cd $1
    go mod init github.com/Zate/go-templates/$1
    go mod tidy
}


# check if we have 0 or 1 arguments

if [ $# -eq 0 ]
then
    tidy
    exit 0
fi

if [ $# -eq 1 ]
then
    new $1
    exit 0
fi