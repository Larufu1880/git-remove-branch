package cmd

import (
	"bytes"
	"os/exec"
)

type GitRepository struct {
	gitBranches []gitBranch
}

type gitBranch struct {
	name      []byte
	removeFlg bool
}

func (gr *GitRepository) DeleteBranches() error {
	for _, gitbranch := range gr.gitBranches {
		if gitbranch.removeFlg {
			if err := exec.Command("git", "branch", "-D", string(gitbranch.name)).Run(); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetBranches display all branches in this repository
func (gr *GitRepository) SetBranches() error {
	gr.gitBranches = []gitBranch{}
	result, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return err
	}
	current, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return err
	}
	current = bytes.Split(current, []byte("\n"))[0]
	branches := bytes.Split(result, []byte("\n"))
	for _, branch := range branches {
		if bytes.Equal(current, branch) {
			continue
		}
		gr.pushBranch(branch)
	}
	return nil
}

func (gr *GitRepository) GetBranchesName() []string {
	var result []string
	for _, gitbranch := range gr.gitBranches {
		result = append(result, string(gitbranch.name))
	}
	return result
}

func (gr *GitRepository) pushBranch(name []byte) {
	if len(name) == 0 {
		return
	}
	gr.gitBranches = append(gr.gitBranches, gitBranch{name, false})
}

func (gr *GitRepository) SetDeleteFlg(index int) {
	if index >= len(gr.gitBranches) {
		return
	}
	gr.gitBranches[index].removeFlg = !gr.gitBranches[index].removeFlg
}

func (gr *GitRepository) GetDeleteFlg() []bool {
	var result []bool
	for _, gitbranch := range gr.gitBranches {
		result = append(result, gitbranch.removeFlg)
	}
	return result
}
