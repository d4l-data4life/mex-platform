package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"

	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	sharedJobs "github.com/d4l-data4life/mex/mex/shared/jobs"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/telemetry"

	pbConfig "github.com/d4l-data4life/mex/mex/services/config/endpoints/config/pb"
)

type RepoCred struct {
	RepoName string
	Token    string
}

type Service struct {
	ServiceTag string
	Log        L.Logger

	Redis *redis.Client

	BroadcastTopicName string

	RepoName          string
	DefaultBranchName string
	EnvPath           string
	UpdateTimeout     time.Duration

	mu              sync.RWMutex
	fs              billy.Filesystem
	currentRepoName string
	currentRepo     *git.Repository
	publicKeys      *gitssh.PublicKeys

	TelemetryService *telemetry.Service
	Jobber           sharedJobs.Jobber

	pbConfig.UnimplementedConfigServer
}

const (
	ConfigResourceName = "config"
	EmptyConfigHash    = "∅"
)

func (svc *Service) InitDeployKey(ctx context.Context, deployKeyPEM []byte) error {
	if len(deployKeyPEM) == 0 {
		return fmt.Errorf("Github deploy key not configured")
	}

	svc.mu.Lock()
	defer svc.mu.Unlock()

	publicKeys, err := gitssh.NewPublicKeys("git", deployKeyPEM, "")
	if err != nil {
		return fmt.Errorf("gitssh: error extracting public key(s): %s", err.Error())
	}
	svc.Log.Info(ctx, L.Messagef("keys: %v", *publicKeys))

	svc.publicKeys = publicKeys
	return nil
}

func (svc *Service) cloneRepo(repoName string) error {
	repoURL := fmt.Sprintf("git@github.com:%s.git", repoName)
	svc.Log.Info(context.Background(), L.Messagef("repo URL: '%s'", repoURL))

	fs := memfs.New()
	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:  repoURL,
		Auth: svc.publicKeys,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	svc.currentRepo = repo
	svc.currentRepoName = repoName
	svc.fs = fs

	return nil
}

func (svc *Service) checkout(remoteName string, refName string) error {
	if svc.currentRepo == nil {
		return fmt.Errorf("repo is nil")
	}

	w, err := svc.currentRepo.Worktree()
	if err != nil {
		return err
	}

	ref := plumbing.NewRemoteReferenceName(remoteName, refName)

	err = svc.currentRepo.Fetch(&git.FetchOptions{
		Auth: svc.publicKeys,
	})
	if err != nil {
		if err != git.NoErrAlreadyUpToDate {
			return E.MakeGRPCStatus(codes.Internal, "git pull failed: "+err.Error(), E.Cause(err)).Err()
		}
		return nil
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: ref,
	})
	if err != nil {
		return E.MakeGRPCStatus(codes.Internal, "git checkout failed: "+err.Error(), E.Cause(err)).Err()
	}

	return nil
}
