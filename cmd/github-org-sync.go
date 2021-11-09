package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type GhOrgSync struct {
	client   *github.Client
	destPath string
	org      string
	wg       sync.WaitGroup
}

func ghClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func fixDestPath(destPath string) string {
	if strings.HasSuffix(destPath, "/") {
		return destPath
	} else {
		return fmt.Sprintf("%s/", destPath)
	}
}

func (gh *GhOrgSync) fetchRepos(ctx context.Context) []*github.Repository {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var repos []*github.Repository
	for {
		repositories, resp, err := gh.client.Repositories.ListByOrg(ctx, gh.org, opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		repos = append(repos, repositories...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return repos
}

func repoExistLocal(repo *github.Repository) bool {
	repoGitDir := fmt.Sprintf("%s/%s/.git", destPath, repo.GetName())
	if _, err := os.Stat(repoGitDir); err != nil {
		return false
	}

	return true
}

func (gh *GhOrgSync) cloneRepo(sem chan struct{}, repo *github.Repository) {
	defer gh.wg.Done()

	sem <- struct{}{}
	defer func() {
		<-sem
	}()

	gitCloneCmd := strings.Fields(fmt.Sprintf("git clone %s %s%s", *repo.CloneURL, gh.destPath, repo.GetName()))
	cmd := exec.Command(gitCloneCmd[0], gitCloneCmd[1:]...)
	fmt.Printf("Cloning repo %s to %s%s\n", repo.GetName(), gh.destPath, repo.GetName())
	err := cmd.Run()
	if err != nil {
		fmt.Printf("***ERROR***: Repo %s failed to clone: %v\n", repo.GetName(), err)
	}
}

func gitDirtyBranch(repoDir string) bool {
	gitDiffIndexCmd := strings.Fields(fmt.Sprintf("git -C %s diff-index --quiet --cached HEAD --", repoDir))
	cmd := exec.Command(gitDiffIndexCmd[0], gitDiffIndexCmd[1:]...)
	err := cmd.Run()
	if err != nil {
		return true
	}

	gitDiffFilesCmd := strings.Fields(fmt.Sprintf("git -C %s diff-files --quiet", repoDir))
	cmd = exec.Command(gitDiffFilesCmd[0], gitDiffFilesCmd[1:]...)
	err = cmd.Run()

	return err != nil
}

func (gh *GhOrgSync) updateLocalRepo(sem chan struct{}, repo *github.Repository) {
	defer gh.wg.Done()

	sem <- struct{}{}
	defer func() {
		<-sem
	}()

	repoPath := fmt.Sprintf("%s%s", gh.destPath, repo.GetName())
	if gitDirtyBranch(fmt.Sprintf("%s%s", gh.destPath, repo.GetName())) {
		fmt.Printf("%s is dirty, so stashing first", repoPath)
		gitStashCmd := strings.Fields(fmt.Sprintf("git -C %s stash push", repoPath))
		cmd := exec.Command(gitStashCmd[0], gitStashCmd[1:]...)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("***ERROR***: Repo %s is dirty but failed to stash: %s\n", repo.GetName(), err)
		}
	}

	defaultBranch := repo.DefaultBranch
	gitCheckoutCmd := strings.Fields(fmt.Sprintf("git -C %s checkout %s", repoPath, *defaultBranch))
	cmd := exec.Command(gitCheckoutCmd[0], gitCheckoutCmd[1:]...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("***ERROR***: Failed to checkout default branch %s for repo %s: %s\n", *defaultBranch, repo.GetName(), err)
	}

	gitUpdateCmd := strings.Fields(fmt.Sprintf("git -C %s pull --rebase", repoPath))
	cmd = exec.Command(gitUpdateCmd[0], gitUpdateCmd[1:]...)
	fmt.Printf("Updating repo %s%s\n", gh.destPath, repo.GetName())
	err = cmd.Run()
	if err != nil {
		fmt.Printf("***ERROR***: Repo %s failed to update: %v\n", repo.GetName(), err)
	}
}

func main(args []string) {
	ctx := context.Background()
	var wg sync.WaitGroup

	gh := GhOrgSync{
		client:   ghClient(ctx),
		destPath: fixDestPath(destPath),
		org:      args[0],
		wg:       wg,
	}

	sem := make(chan struct{}, parallelism)
	repos := gh.fetchRepos(ctx)

	for _, repo := range repos {
		if *repo.Archived && !includeArchived {
			fmt.Printf("Not including %s since it is archived and you asked to not included archived\n", repo.GetName())
			continue
		} else {
			gh.wg.Add(1)
			if repoExistLocal(repo) {
				go gh.updateLocalRepo(sem, repo)
			} else {
				go gh.cloneRepo(sem, repo)
			}
		}
	}
	gh.wg.Wait()
}
