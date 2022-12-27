package events_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/runatlantis/atlantis/server/events"
	"github.com/runatlantis/atlantis/server/events/models"
	"github.com/runatlantis/atlantis/server/logging"
	. "github.com/runatlantis/atlantis/testing"
)

// disableSSLVerification disables ssl verification for the global http client
// and returns a function to be called in a defer that will re-enable it.
func disableSSLVerification() func() {
	orig := http.DefaultTransport.(*http.Transport).TLSClientConfig
	// nolint: gosec
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return func() {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = orig
	}
}

// Test that if we don't have any existing files, we check out the repo.
func TestClone_NoneExisting(t *testing.T) {
	// Initialize the git repo.
	repoDir := initRepo(t)
	expCommit := runCmd(t, repoDir, "git", "rev-parse", "HEAD")

	dataDir := t.TempDir()

	wd := &events.FileWorkspace{
		DataDir:                     dataDir,
		CheckoutMerge:               false,
		TestingOverrideHeadCloneURL: fmt.Sprintf("file://%s", repoDir),
		GpgNoSigningEnabled:         true,
	}

	cloneDir, _, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
	}, "default", ".")
	Ok(t, err)

	// Use rev-parse to verify at correct commit.
	actCommit := runCmd(t, cloneDir, "git", "rev-parse", "HEAD")
	Equals(t, expCommit, actCommit)
}

// Test that if we don't have any existing files, we check out the repo
// successfully when we're using the merge method.
func TestClone_CheckoutMergeNoneExisting(t *testing.T) {
	// Initialize the git repo.
	repoDir := initRepo(t)

	// Add a commit to branch 'branch' that's not on main.
	runCmd(t, repoDir, "git", "checkout", "branch")
	runCmd(t, repoDir, "touch", "branch-file")
	runCmd(t, repoDir, "git", "add", "branch-file")
	runCmd(t, repoDir, "git", "commit", "-m", "branch-commit")
	branchCommit := runCmd(t, repoDir, "git", "rev-parse", "HEAD")

	// Now switch back to main and advance the main branch by another
	// commit.
	runCmd(t, repoDir, "git", "checkout", "main")
	runCmd(t, repoDir, "touch", "main-file")
	runCmd(t, repoDir, "git", "add", "main-file")
	runCmd(t, repoDir, "git", "commit", "-m", "main-commit")
	mainCommit := runCmd(t, repoDir, "git", "rev-parse", "HEAD")

	// Finally, perform a merge in another branch ourselves, just so we know
	// what the final state of the repo should be.
	runCmd(t, repoDir, "git", "checkout", "-b", "mergetest")
	runCmd(t, repoDir, "git", "merge", "-m", "atlantis-merge", "branch")
	expLsOutput := runCmd(t, repoDir, "ls")

	dataDir := t.TempDir()

	overrideURL := fmt.Sprintf("file://%s", repoDir)
	wd := &events.FileWorkspace{
		DataDir:                     dataDir,
		CheckoutMerge:               true,
		TestingOverrideHeadCloneURL: overrideURL,
		TestingOverrideBaseCloneURL: overrideURL,
		GpgNoSigningEnabled:         true,
	}

	cloneDir, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
		BaseBranch: "main",
	}, "default", ".")
	Ok(t, err)
	Equals(t, false, hasDiverged)

	// Check the commits.
	actBaseCommit := runCmd(t, cloneDir, "git", "rev-parse", "HEAD~1")
	actHeadCommit := runCmd(t, cloneDir, "git", "rev-parse", "HEAD^2")
	Equals(t, mainCommit, actBaseCommit)
	Equals(t, branchCommit, actHeadCommit)

	// Use ls to verify the repo looks good.
	actLsOutput := runCmd(t, cloneDir, "ls")
	Equals(t, expLsOutput, actLsOutput)
}

// Test that if we're using the merge method and the repo is already cloned at
// the right commit, then we don't reclone.
func TestClone_CheckoutMergeNoReclone(t *testing.T) {
	// Initialize the git repo.
	repoDir := initRepo(t)

	// Add a commit to branch 'branch' that's not on main.
	runCmd(t, repoDir, "git", "checkout", "branch")
	runCmd(t, repoDir, "touch", "branch-file")
	runCmd(t, repoDir, "git", "add", "branch-file")
	runCmd(t, repoDir, "git", "commit", "-m", "branch-commit")

	// Now switch back to main and advance the main branch by another commit.
	runCmd(t, repoDir, "git", "checkout", "main")
	runCmd(t, repoDir, "touch", "main-file")
	runCmd(t, repoDir, "git", "add", "main-file")
	runCmd(t, repoDir, "git", "commit", "-m", "main-commit")

	// Run the clone for the first time.
	dataDir := t.TempDir()
	overrideURL := fmt.Sprintf("file://%s", repoDir)
	wd := &events.FileWorkspace{
		DataDir:                     dataDir,
		CheckoutMerge:               true,
		TestingOverrideHeadCloneURL: overrideURL,
		TestingOverrideBaseCloneURL: overrideURL,
		GpgNoSigningEnabled:         true,
	}

	_, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
		BaseBranch: "main",
	}, "default", ".")
	Ok(t, err)
	Equals(t, false, hasDiverged)

	// Create a file that we can use to check if the repo was recloned.
	runCmd(t, dataDir, "touch", "repos/0/default/FY======/proof")

	// Now run the clone again.
	cloneDir, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
		BaseBranch: "main",
	}, "default", ".")
	Ok(t, err)
	Equals(t, false, hasDiverged)

	// Check that our proof file is still there, proving that we didn't reclone.
	_, err = os.Stat(filepath.Join(cloneDir, "proof"))
	Ok(t, err)
}

// Same as TestClone_CheckoutMergeNoReclone however the branch that gets
// merged is a fast-forward merge. See #584.
func TestClone_CheckoutMergeNoRecloneFastForward(t *testing.T) {
	// Initialize the git repo.
	repoDir := initRepo(t)

	// Add a commit to branch 'branch' that's not on main.
	// This will result in a fast-forwardable merge.
	runCmd(t, repoDir, "git", "checkout", "branch")
	runCmd(t, repoDir, "touch", "branch-file")
	runCmd(t, repoDir, "git", "add", "branch-file")
	runCmd(t, repoDir, "git", "commit", "-m", "branch-commit")

	// Run the clone for the first time.
	dataDir := t.TempDir()
	overrideURL := fmt.Sprintf("file://%s", repoDir)
	wd := &events.FileWorkspace{
		DataDir:                     dataDir,
		CheckoutMerge:               true,
		TestingOverrideHeadCloneURL: overrideURL,
		TestingOverrideBaseCloneURL: overrideURL,
		GpgNoSigningEnabled:         true,
	}

	_, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
		BaseBranch: "main",
	}, "default", ".")
	Ok(t, err)
	Equals(t, false, hasDiverged)

	// Create a file that we can use to check if the repo was recloned.
	runCmd(t, dataDir, "touch", "repos/0/default/FY======/proof")

	// Now run the clone again.
	cloneDir, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
		BaseBranch: "main",
	}, "default", ".")
	Ok(t, err)
	Equals(t, false, hasDiverged)

	// Check that our proof file is still there, proving that we didn't reclone.
	_, err = os.Stat(filepath.Join(cloneDir, "proof"))
	Ok(t, err)
}

// Test that if there's a conflict when merging we return a good error.
func TestClone_CheckoutMergeConflict(t *testing.T) {
	// Initialize the git repo.
	repoDir := initRepo(t)

	// Add a commit to branch 'branch' that's not on main.
	runCmd(t, repoDir, "git", "checkout", "branch")
	runCmd(t, repoDir, "sh", "-c", "echo hi >> file")
	runCmd(t, repoDir, "git", "add", "file")
	runCmd(t, repoDir, "git", "commit", "-m", "branch-commit")

	// Add a new commit to main that will cause a conflict if branch was
	// merged.
	runCmd(t, repoDir, "git", "checkout", "main")
	runCmd(t, repoDir, "sh", "-c", "echo conflict >> file")
	runCmd(t, repoDir, "git", "add", "file")
	runCmd(t, repoDir, "git", "commit", "-m", "commit")

	// We're set up, now trigger the Atlantis clone.
	dataDir := t.TempDir()
	overrideURL := fmt.Sprintf("file://%s", repoDir)
	wd := &events.FileWorkspace{
		DataDir:                     dataDir,
		CheckoutMerge:               true,
		TestingOverrideHeadCloneURL: overrideURL,
		TestingOverrideBaseCloneURL: overrideURL,
		GpgNoSigningEnabled:         true,
	}

	_, _, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
		BaseBranch: "main",
	}, "default", ".")

	ErrContains(t, "running git merge -q --no-ff -m atlantis-merge FETCH_HEAD", err)
	ErrContains(t, "Auto-merging file", err)
	ErrContains(t, "CONFLICT (add/add)", err)
	ErrContains(t, "Merge conflict in file", err)
	ErrContains(t, "Automatic merge failed; fix conflicts and then commit the result.", err)
	ErrContains(t, "exit status 1", err)
}

// Test that if the repo is already cloned and is at the right commit, we
// don't reclone.
func TestClone_NoReclone(t *testing.T) {
	repoDir := initRepo(t)
	dataDir := t.TempDir()

	runCmd(t, dataDir, "mkdir", "-p", "repos/0/default")
	runCmd(t, dataDir, "mv", repoDir, "repos/0/default/FY======")
	// Create a file that we can use later to check if the repo was recloned.
	runCmd(t, dataDir, "touch", "repos/0/default/FY======/proof")

	wd := &events.FileWorkspace{
		DataDir:                     dataDir,
		CheckoutMerge:               false,
		TestingOverrideHeadCloneURL: fmt.Sprintf("file://%s", repoDir),
		GpgNoSigningEnabled:         true,
	}
	cloneDir, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
	}, "default", ".")
	Ok(t, err)
	Equals(t, false, hasDiverged)

	// Check that our proof file is still there.
	_, err = os.Stat(filepath.Join(cloneDir, "proof"))
	Ok(t, err)
}

// Test that if the repo is already cloned but is at the wrong commit, we
// reclone.
func TestClone_RecloneWrongCommit(t *testing.T) {
	repoDir := initRepo(t)
	dataDir := t.TempDir()

	// Copy the repo to our data dir.
	runCmd(t, dataDir, "mkdir", "-p", "repos/0/default")
	runCmd(t, dataDir, "cp", "-R", repoDir, "repos/0/default/FY======")

	// Now add a commit to the repo, so the one in the data dir is out of date.
	runCmd(t, repoDir, "git", "checkout", "branch")
	runCmd(t, repoDir, "touch", "newfile")
	runCmd(t, repoDir, "git", "add", "newfile")
	runCmd(t, repoDir, "git", "commit", "-m", "newfile")
	expCommit := runCmd(t, repoDir, "git", "rev-parse", "HEAD")

	wd := &events.FileWorkspace{
		DataDir:                     dataDir,
		CheckoutMerge:               false,
		TestingOverrideHeadCloneURL: fmt.Sprintf("file://%s", repoDir),
		GpgNoSigningEnabled:         true,
	}
	cloneDir, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "branch",
		HeadCommit: expCommit,
	}, "default", ".")
	Ok(t, err)
	Equals(t, false, hasDiverged)

	// Use rev-parse to verify at correct commit.
	actCommit := runCmd(t, cloneDir, "git", "rev-parse", "HEAD")
	Equals(t, expCommit, actCommit)
}

// Test that if the branch we're merging into has diverged and we're using
// checkout-strategy=merge, we warn the user (see #804).
func TestClone_mainHasDiverged(t *testing.T) {
	// Initialize the git repo.
	repoDir := initRepo(t)

	// Simulate first PR.
	runCmd(t, repoDir, "git", "checkout", "-b", "first-pr")
	runCmd(t, repoDir, "touch", "file1")
	runCmd(t, repoDir, "git", "add", "file1")
	runCmd(t, repoDir, "git", "commit", "-m", "file1")

	// Atlantis checkout first PR.
	firstPRDir := repoDir + "/first-pr"
	runCmd(t, repoDir, "mkdir", "-p", "first-pr")
	runCmd(t, firstPRDir, "git", "clone", "--branch", "main", "--single-branch", repoDir, ".")
	runCmd(t, firstPRDir, "git", "remote", "add", "head", repoDir)
	runCmd(t, firstPRDir, "git", "fetch", "head", "+refs/heads/first-pr")
	runCmd(t, firstPRDir, "git", "config", "--local", "user.email", "atlantisbot@runatlantis.io")
	runCmd(t, firstPRDir, "git", "config", "--local", "user.name", "atlantisbot")
	runCmd(t, firstPRDir, "git", "config", "--local", "commit.gpgsign", "false")
	runCmd(t, firstPRDir, "git", "merge", "-q", "--no-ff", "-m", "atlantis-merge", "FETCH_HEAD")

	// Simulate second PR.
	runCmd(t, repoDir, "git", "checkout", "main")
	runCmd(t, repoDir, "git", "checkout", "-b", "second-pr")
	runCmd(t, repoDir, "touch", "file2")
	runCmd(t, repoDir, "git", "add", "file2")
	runCmd(t, repoDir, "git", "commit", "-m", "file2")

	// Atlantis checkout second PR.
	secondPRDir := repoDir + "/second-pr"
	runCmd(t, repoDir, "mkdir", "-p", "second-pr")
	runCmd(t, secondPRDir, "git", "clone", "--branch", "main", "--single-branch", repoDir, ".")
	runCmd(t, secondPRDir, "git", "remote", "add", "head", repoDir)
	runCmd(t, secondPRDir, "git", "fetch", "head", "+refs/heads/second-pr")
	runCmd(t, secondPRDir, "git", "config", "--local", "user.email", "atlantisbot@runatlantis.io")
	runCmd(t, secondPRDir, "git", "config", "--local", "user.name", "atlantisbot")
	runCmd(t, secondPRDir, "git", "config", "--local", "commit.gpgsign", "false")
	runCmd(t, secondPRDir, "git", "merge", "-q", "--no-ff", "-m", "atlantis-merge", "FETCH_HEAD")

	// Merge first PR
	runCmd(t, repoDir, "git", "checkout", "main")
	runCmd(t, repoDir, "git", "merge", "first-pr")

	// Copy the second-pr repo to our data dir which has diverged remote main
	runCmd(t, repoDir, "mkdir", "-p", "repos/0/default")
	runCmd(t, repoDir, "cp", "-R", secondPRDir, "repos/0/default/FY======")

	// Run the clone.
	wd := &events.FileWorkspace{
		DataDir:             repoDir,
		CheckoutMerge:       true,
		GpgNoSigningEnabled: true,
	}
	_, hasDiverged, err := wd.Clone(logging.NewNoopLogger(t), models.Repo{CloneURL: repoDir}, models.PullRequest{
		BaseRepo:   models.Repo{CloneURL: repoDir},
		HeadBranch: "second-pr",
		BaseBranch: "main",
	}, "default", ".")
	Ok(t, err)
	Equals(t, hasDiverged, true)

	// Run it again but without the checkout merge strategy. It should return
	// false.
	wd.CheckoutMerge = false
	_, hasDiverged, err = wd.Clone(logging.NewNoopLogger(t), models.Repo{}, models.PullRequest{
		BaseRepo:   models.Repo{},
		HeadBranch: "second-pr",
		BaseBranch: "main",
	}, "default", ".")
	Ok(t, err)
	Equals(t, hasDiverged, false)
}

func TestHasDiverged_mainHasDiverged(t *testing.T) {
	// Initialize the git repo.
	repoDir := initRepo(t)

	// Simulate first PR.
	runCmd(t, repoDir, "git", "checkout", "-b", "first-pr")
	runCmd(t, repoDir, "touch", "file1")
	runCmd(t, repoDir, "git", "add", "file1")
	runCmd(t, repoDir, "git", "commit", "-m", "file1")

	// Atlantis checkout first PR.
	firstPRDir := repoDir + "/first-pr"
	runCmd(t, repoDir, "mkdir", "-p", "first-pr")
	runCmd(t, firstPRDir, "git", "clone", "--branch", "main", "--single-branch", repoDir, ".")
	runCmd(t, firstPRDir, "git", "remote", "add", "head", repoDir)
	runCmd(t, firstPRDir, "git", "fetch", "head", "+refs/heads/first-pr")
	runCmd(t, firstPRDir, "git", "config", "--local", "user.email", "atlantisbot@runatlantis.io")
	runCmd(t, firstPRDir, "git", "config", "--local", "user.name", "atlantisbot")
	runCmd(t, firstPRDir, "git", "config", "--local", "commit.gpgsign", "false")
	runCmd(t, firstPRDir, "git", "merge", "-q", "--no-ff", "-m", "atlantis-merge", "FETCH_HEAD")

	// Simulate second PR.
	runCmd(t, repoDir, "git", "checkout", "main")
	runCmd(t, repoDir, "git", "checkout", "-b", "second-pr")
	runCmd(t, repoDir, "touch", "file2")
	runCmd(t, repoDir, "git", "add", "file2")
	runCmd(t, repoDir, "git", "commit", "-m", "file2")

	// Atlantis checkout second PR.
	secondPRDir := repoDir + "/second-pr"
	runCmd(t, repoDir, "mkdir", "-p", "second-pr")
	runCmd(t, secondPRDir, "git", "clone", "--branch", "main", "--single-branch", repoDir, ".")
	runCmd(t, secondPRDir, "git", "remote", "add", "head", repoDir)
	runCmd(t, secondPRDir, "git", "fetch", "head", "+refs/heads/second-pr")
	runCmd(t, secondPRDir, "git", "config", "--local", "user.email", "atlantisbot@runatlantis.io")
	runCmd(t, secondPRDir, "git", "config", "--local", "user.name", "atlantisbot")
	runCmd(t, secondPRDir, "git", "config", "--local", "commit.gpgsign", "false")
	runCmd(t, secondPRDir, "git", "merge", "-q", "--no-ff", "-m", "atlantis-merge", "FETCH_HEAD")

	// Merge first PR
	runCmd(t, repoDir, "git", "checkout", "main")
	runCmd(t, repoDir, "git", "merge", "first-pr")

	// Copy the second-pr repo to our data dir which has diverged remote main
	runCmd(t, repoDir, "mkdir", "-p", "repos/0/default")
	runCmd(t, repoDir, "cp", "-R", secondPRDir, "repos/0/default/FY======")

	// "git", "remote", "set-url", "origin", p.BaseRepo.CloneURL,
	runCmd(t, repoDir+"/repos/0/default/FY======", "git", "remote", "update")

	// Run the clone.
	wd := &events.FileWorkspace{
		DataDir:             repoDir,
		CheckoutMerge:       true,
		GpgNoSigningEnabled: true,
	}
	hasDiverged := wd.HasDiverged(logging.NewNoopLogger(t), repoDir+"/repos/0/default/FY======")
	Equals(t, hasDiverged, true)

	// Run it again but without the checkout merge strategy. It should return
	// false.
	wd.CheckoutMerge = false
	hasDiverged = wd.HasDiverged(logging.NewNoopLogger(t), repoDir+"/repos/0/default/FY======")
	Equals(t, hasDiverged, false)
}

func initRepo(t *testing.T) string {
	repoDir := t.TempDir()
	runCmd(t, repoDir, "git", "init", "--initial-branch=main")
	runCmd(t, repoDir, "touch", ".gitkeep")
	runCmd(t, repoDir, "git", "add", ".gitkeep")
	runCmd(t, repoDir, "git", "config", "--local", "user.email", "atlantisbot@runatlantis.io")
	runCmd(t, repoDir, "git", "config", "--local", "user.name", "atlantisbot")
	runCmd(t, repoDir, "git", "config", "--local", "commit.gpgsign", "false")
	runCmd(t, repoDir, "git", "commit", "-m", "initial commit")
	runCmd(t, repoDir, "git", "branch", "branch")
	return repoDir
}
