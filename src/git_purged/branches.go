// This file is a part of 'git-purged' tool, http://thekondor.net

package git_purged

type Branches struct {
    git Git
}

const GoneTrackName = "[gone]"

func filter(branches []GitBranch, pred func(GitBranch) bool) []string {
    filtered := []string{}

    for _, branch := range branches {
        if pred(branch) {
            filtered = append(filtered, branch.Name)
        }
    }

    return filtered
}

func (self Branches) list(pred func(GitBranch) bool) ([]string, error) {
    allBranches, err := self.git.ListLocalBranches()
    if nil != err {
        return nil, err
    }

    return filter(allBranches, pred), nil
}

func (self Branches) ListGone() ([]string, error) {
    return self.list(func(branch GitBranch) bool {
        return GoneTrackName == branch.Track
    })
}

func (self Branches) ListNonGone() ([]string, error) {
    return self.list(func(branch GitBranch) bool {
        return GoneTrackName != branch.Track
    })
}

func NewBranches(git Git) Branches {
    return Branches{ git : git }
}

