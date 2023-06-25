package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func ListChangedFiles(srcDir string, relBranch string) ([]string, error) {
	repo, err := git.PlainOpen(srcDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the srcDir git data: %w", err)
	}
	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve HEAD ref: %w", err)
	}
	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get the commit object for HEAD ref: %w", err)
	}
	mainRef, err := repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", relBranch)), true)
	if err != nil {
		return nil, fmt.Errorf("failed to get the relative branch ref: %w", err)
	}
	mainHead, err := repo.CommitObject(mainRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get the commit object for relative branch: %w", err)
	}
	diff, err := commit.Patch(mainHead)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff between HEAD and relative branch: %w", err)
	}
	var files []string
	for _, stat := range diff.Stats() {
		files = append(files, stat.Name)
	}
	return files, nil
}
