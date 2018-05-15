// This file is a part of 'git-purged' tool, http://thekondor.net

package git_purged

import (
    "fmt"
    "os/exec"
    "bufio"
    "strings"
)

type GitExternalCommand struct {
}

func collectOutput(scanner *bufio.Scanner) []string {
    var output []string
    for scanner.Scan() {
        output = append(output, scanner.Text())
    }

    return output
}

func (self GitExternalCommand) Run(args ...string) ([]string, error) {
    cmd := exec.Command("git", args...)

    stdoutPipe, err := cmd.StdoutPipe()
    if nil != err {
        return nil, fmt.Errorf("Git: failed to obtain stdout pipe for stdout: %v", err)
    }

    stderrPipe, err := cmd.StderrPipe()
    if nil != err {
        return nil, fmt.Errorf("Git: failed to obtain stderr pipe for stderr: %v", err)
    }

    cmd.Start()

    stdout := collectOutput(bufio.NewScanner(stdoutPipe))
    stderr := collectOutput(bufio.NewScanner(stderrPipe))

    cmd.Wait()

    if len(stderr) > 0 {
        return nil, fmt.Errorf("Git: '%s' failed: '%s'", strings.Join(args, " "), strings.Join(stderr, "\n"))
    }

    return stdout, nil
}

func NewGitExternalCommand() GitExternalCommand {
    return GitExternalCommand{}
}

