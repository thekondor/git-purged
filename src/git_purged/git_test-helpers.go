// This file is a part of 'git-purged' tool, http://thekondor.net

package git_purged

import (
    "github.com/stretchr/testify/mock"
)

func arrayOfStrings(args mock.Arguments, index int) []string {
    arg := args.Get(0)
    arrayOfString, ok := arg.([]string)
    if !ok {
        panic("Cannot cast to array of strings")
    }

    return arrayOfString
}

var EmptyStdOut = []string{}
var EmptyError error = nil

func withArgs(args ...string) []string {
    return args
}

func stdOut(args ...string) []string {
    return args
}

func anyArgument(interface{}) bool {
    return true
}

