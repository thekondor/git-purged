// This file is a part of 'git-purged' tool, https://thekondor.net

package main

import (
    "git_purged"
    "os"
    "fmt"
    "flag"
    logger "log"
)

type App struct {
    git *git_purged.Git
    branches git_purged.Branches
    log *logger.Logger
}

func NewApp(git *git_purged.Git, log *logger.Logger) App {
    return App{
        git : git,
        branches : git_purged.NewBranches(*git),
        log : log,
    }
}

func (self App) PrintPurgedBranches(pruneBefore bool) {
    if pruneBefore {
        self.pruneRemoteOrigin()
    }

    goneBranchNames, err := self.branches.ListGone()
    if nil != err {
        self.log.Fatalf("Failed to list purged branches: %v\n", err)
    }

    self.list(goneBranchNames)
}

func (self App) PrintAliveBranches(pruneBefore bool) {
    if pruneBefore {
        self.pruneRemoteOrigin()
    }

    nonGoneBranchNames, err := self.branches.ListNonGone()
    if nil != err {
        self.log.Fatalf("Failed to list non-purged branches: %v\n", err)
    }

    self.list(nonGoneBranchNames)
}

func (self App) list(branchNames []string) {
    if 0 == len(branchNames) {
        self.log.Printf("No branches to show. Didn't you forget to call 'git fetch' before?")
        return
    }

    for _, name := range branchNames {
        fmt.Println(name)
    }
}

func (self App) pruneRemoteOrigin() {
    if err := self.git.PruneRemoteOrigin(); nil != err {
        self.log.Fatalf("Failed to prune remote origin: %v\n", err)
    }
}

func NewGit(log *logger.Logger) *git_purged.Git{
    git, err := git_purged.NewGit(git_purged.NewGitExternalCommand())
    if nil == err {
        return git
    }

    if git_purged.GitNotAvailableErr == err {
        log.Fatal("Error: 'git' command is not available or misconfigured");
    } else {
        log.Fatalf("Error: 'git' command cannot be used (%v)\n", err)
    }

    return nil
}

func main() {
    log := logger.New(os.Stdout, "[git-purged] ", 0)

    inverseFlag := flag.Bool("inverse", false, "Inverse output by printing alive branches only. Optional.")
    skipPruneFlag := flag.Bool("skip-prune", false, "Skip prunning of remote's origin to calculate purged branches before listing. Optional.")
    helpFlag := flag.Bool("help", false, "Show this help.")

    flag.Parse()
    if *helpFlag {
        fmt.Printf("git-purged - subcommand to list purged (already removed on a remote master) branches.\n")
        fmt.Printf("             build: %s#%s (%s)\n\n", BuildDate, BuildGitCommit, BuildGitOrigin)
        flag.PrintDefaults()
        fmt.Println()
        log.Fatalf("No action requested")
    }

    git := NewGit(log)
    if !git.IsRepository() {
        log.Fatalf("Current directory is not a valid working git tree")
    }

    app := NewApp(git, log)
    if *inverseFlag {
        app.PrintAliveBranches(*skipPruneFlag)
    } else {
        app.PrintPurgedBranches(*skipPruneFlag)
    }
}

