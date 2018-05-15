// This file is a part of 'git-purged' tool, http://thekondor.net

package git_purged

import (
    "fmt"
    "strings"
)

type Git struct {
    binary string
    command ExternalCommand
}

type GitBranch struct {
    Track string
    Name string
}

var (
    GitNotAvailableErr = fmt.Errorf("GIT: 'git' is not available in PATH or misconfigured")
)

type ExternalCommand interface {
    Run(args ...string) (stdout []string, err error)
}

func NewGit(gitCommand ExternalCommand) (*Git, error) {
    _, err := gitCommand.Run("--version")

    if nil != err {
        return nil, GitNotAvailableErr
    }

    return &Git{ command : gitCommand }, nil
}

func (self Git) IsRepository() bool {
    stdout, err := self.command.Run("rev-parse", "--is-inside-work-tree")
    if nil != err {
        return false
    }

    return len(stdout) > 0 && "true" == stdout[0]
}

func (self Git) PruneRemoteOrigin() error {
    _, err := self.command.Run("remote", "prune", "origin")
    if nil != err {
        return fmt.Errorf("Git: failed to prune remote origin's references: %v", err)
    }

    return nil
}

func (self Git) ListLocalBranches() ([]GitBranch, error) {
    stdout, err := self.command.Run("branch", "-vv", "--format=%(upstream:track):%(refname:lstrip=2)")
    if nil != err {
        return nil, fmt.Errorf("Git: failed to execute branch listing: %v", err)
    }

    var parsed []GitBranch
    for _, line := range stdout {
        parsedLine := strings.Split(line, ":")
        if 2 != len(parsedLine) {
            return nil, fmt.Errorf("Git: invalid local branches format: %s", line)
        }

        parsed = append(parsed, GitBranch{ parsedLine[0], parsedLine[1] })
    }

    return parsed, nil
}

