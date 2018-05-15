// This file is a part of 'git-purged' tool, http://thekondor.net

package git_purged

import (
    "testing"
    "errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestGit_CreationFailsOnAvailabilityCheck(t *testing.T) {
    mockGit := new(GitExternalCommandMock)
    mockGit.On("Run", withArgs("--version")).
        Return(
            EmptyStdOut,
            errors.New("Command not found"),
        )

    sut, err := NewGit(mockGit)

    require.Nil(t, sut, "Nil is returned on error")
    assert.NotNil(t, err, "Error is returned on git non available")
    assert.Contains(t, err.Error(), "not available")
}

func mockGitCommand() *GitExternalCommandMock {
    mockGit := new(GitExternalCommandMock)
    mockGit.On("Run", mock.MatchedBy(anyArgument)).
        Return(
            EmptyStdOut, EmptyError,
        ).Once()

    return mockGit
}

func prepareSut(t *testing.T) (*Git, *GitExternalCommandMock) {
    mockGit := mockGitCommand()
    sut, err := NewGit(mockGit)

    assert.Nil(t, err)
    assert.NotNil(t, sut)

    return sut, mockGit
}

func TestGit_CreationSucceededOnAvailabiltyCheck(t *testing.T) {
    mockGit := new(GitExternalCommandMock)
    mockGit.On("Run", withArgs("--version")).
        Return(
            EmptyStdOut, EmptyError,
        )

    sut, err := NewGit(mockGit)

    assert.NotNil(t, sut, "Valid object is returned")
    require.Nil(t, err)
}

func TestGit_ErrorIsPassedOver_OnFailedPrune(t *testing.T) {
    sut, mockGit := prepareSut(t)

    mockGit.On("Run", withArgs("remote", "prune", "origin")).
        Return(
            EmptyStdOut, errors.New("test:prune failed"),
        ).Once()

    err := sut.PruneRemoteOrigin()

    require.NotNil(t, err, "Error is returned on prune remote origin")
    require.Contains(t, err.Error(), "test:prune failed")
}

func TestGit_ErrorIsPassedOver_OnFailedBranchListing(t *testing.T) {
    sut, mockGit := prepareSut(t)

    mockGit.On("Run", withArgs("branch", "-vv", "--format=%(upstream:track):%(refname:lstrip=2)")).
        Return(
            EmptyStdOut, errors.New("test:branch listing failed"),
        ).Once()

    _, err := sut.ListLocalBranches()

    require.NotNil(t, err, "Error is returned on listing local branches")
    require.Contains(t, err.Error(), "test:branch listing failed")
}

func TestGit_BranchesAreReturned_OnBranchListing(t *testing.T) {
    sut, mockGit := prepareSut(t)

    mockGit.On("Run", mock.MatchedBy(anyArgument)).
        Return(
            stdOut(
                ":refs/heads/branch-1",
                "[gone]:refs/heads/branch-2",
            ), EmptyError,
        ).Once()

    branches, _ := sut.ListLocalBranches()
    require.Len(t, branches, 2, "Amount of listed branches is recognized")
}

func TestGit_GoneBranchesAreRecognized_OnBranchListing(t *testing.T) {
    sut, mockGit := prepareSut(t)

    mockGit.On("Run", mock.MatchedBy(anyArgument)).
        Return(
            stdOut(
                ":refs/heads/branch-1",
                "[gone]:refs/heads/branch-2",
            ), EmptyError,
        ).Once()

    branches, _ := sut.ListLocalBranches()
    assert.NotEmpty(t, branches)
    require.Contains(t, branches, GitBranch{"[gone]", "refs/heads/branch-2"}, "Gone branch is recognized")
}

