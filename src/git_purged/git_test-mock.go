// This file is a part of 'git-purged' tool, http://thekondor.net

package git_purged

import (
    "github.com/stretchr/testify/mock"
)

type GitExternalCommandMock struct {
    mock.Mock
}

func (self *GitExternalCommandMock) Run(args ...string) ([]string, error) {
    mockArgs := self.Called(args)
    return arrayOfStrings(mockArgs, 0), mockArgs.Error(1)
}

