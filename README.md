# git-purged
git-purged - a git sub command to list branches already removed on remote origin

## Overview

A tiny tool to list branches which are already purged from a remote origin. A cross-platform and easy-to-use alternative for grep'ing `git`'s output. The main purpose of the `git-purged` is to provide a convinient way for scripting deletion of obsolete git branches.

## Build

Since the tool was intentionally written with using no external dependencies, the following should be enough to a get a self-contained binary:

    make build
    
will produce `git-purged` (`git-purged.exe` for MS Windows) binary in `build/` directory.

## Usage

The only requirement is `git` command availability through `${PATH}` (`%PATH%`) global environment variable.

### List purged

Run with not arguments. Default behavior.

### List non-purged

    git-purged --inverse

# License

MIT

