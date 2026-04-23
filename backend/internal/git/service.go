package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Service 提供 Git 仓库操作服务
type Service struct {
	repoBasePath string
}

// NewService 创建 Git 服务实例
func NewService(basePath string) *Service {
	return &Service{repoBasePath: basePath}
}

// getRepoPath 获取仓库存储路径
func (s *Service) getRepoPath(owner, name string) string {
	return filepath.Join(s.repoBasePath, owner, name+".git")
}

// CreateRepo 创建新的 Git 仓库
func (s *Service) CreateRepo(ctx context.Context, owner, name string, initWithReadme bool) error {
	repoPath := s.getRepoPath(owner, name)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(repoPath), 0755); err != nil {
		return fmt.Errorf("创建仓库目录失败: %w", err)
	}

	// 创建仓库
	_, err := git.PlainInit(repoPath, true) // true 表示 bare 仓库
	if err != nil {
		return fmt.Errorf("初始化仓库失败: %w", err)
	}

	// 如果需要初始化 README，我们需要创建一个非 bare 仓库并推送到 bare 仓库
	if initWithReadme {
		if err := s.initWithReadme(repoPath, owner, name); err != nil {
			return err
		}
	}

	return nil
}

// initWithReadme 初始化带有 README 的仓库
func (s *Service) initWithReadme(bareRepoPath, owner, name string) error {
	// 创建临时目录用于初始化
	tempDir, err := os.MkdirTemp("", "laima-repo-init-*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// 初始化普通仓库
	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		return fmt.Errorf("初始化临时仓库失败: %w", err)
	}

	// 创建 README 文件
	readmePath := filepath.Join(tempDir, "README.md")
	readmeContent := fmt.Sprintf("# %s\n\nThis is a Laima repository.", name)
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("创建 README 失败: %w", err)
	}

	// 获取工作区
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("获取工作区失败: %w", err)
	}

	// 添加文件
	_, err = w.Add("README.md")
	if err != nil {
		return fmt.Errorf("添加文件失败: %w", err)
	}

	// 提交
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  owner,
			Email: fmt.Sprintf("%s@laima.local", owner),
		},
	})
	if err != nil {
		return fmt.Errorf("提交失败: %w", err)
	}

	// 添加远程仓库
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{bareRepoPath},
	})
	if err != nil {
		return fmt.Errorf("添加远程仓库失败: %w", err)
	}

	// 推送到 bare 仓库
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"refs/heads/*:refs/heads/*"},
	})
	if err != nil {
		return fmt.Errorf("推送失败: %w", err)
	}

	return nil
}

// DeleteRepo 删除 Git 仓库
func (s *Service) DeleteRepo(ctx context.Context, owner, name string) error {
	repoPath := s.getRepoPath(owner, name)
	if err := os.RemoveAll(repoPath); err != nil {
		return fmt.Errorf("删除仓库失败: %w", err)
	}
	return nil
}

// GetRepo 获取 Git 仓库实例
func (s *Service) GetRepo(owner, name string) (*git.Repository, error) {
	repoPath := s.getRepoPath(owner, name)
	return git.PlainOpen(repoPath)
}

// GetCommit 获取提交信息
func (s *Service) GetCommit(owner, name, commitHash string) (*object.Commit, error) {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}

	hash := plumbing.NewHash(commitHash)
	return repo.CommitObject(hash)
}

// GetFileContent 获取文件内容
func (s *Service) GetFileContent(owner, name, ref, path string) (string, error) {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return "", err
	}

	// 获取引用
	refObj, err := repo.Reference(plumbing.ReferenceName(ref), true)
	if err != nil {
		return "", err
	}

	// 获取提交
	commit, err := repo.CommitObject(refObj.Hash())
	if err != nil {
		return "", err
	}

	// 获取文件
	file, err := commit.File(path)
	if err != nil {
		return "", err
	}

	return file.Contents()
}

// ListFiles 列出目录中的文件
func (s *Service) ListFiles(owner, name, ref, path string) ([]string, error) {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}

	refObj, err := repo.Reference(plumbing.ReferenceName(ref), true)
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(refObj.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	var files []string
	if path == "" || path == "/" {
		// 根目录
		for _, entry := range tree.Entries {
			files = append(files, entry.Name)
		}
	} else {
		// 子目录
		subTree, err := tree.Tree(path)
		if err != nil {
			return nil, err
		}
		for _, entry := range subTree.Entries {
			files = append(files, entry.Name)
		}
	}

	return files, nil
}

// CreateBranch 创建新分支
func (s *Service) CreateBranch(owner, name, branchName, sourceRef string) error {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return err
	}

	ref, err := repo.Reference(plumbing.ReferenceName(sourceRef), true)
	if err != nil {
		return err
	}

	branchRef := plumbing.NewHashReference(
		plumbing.ReferenceName("refs/heads/"+branchName), ref.Hash())
	return repo.Storer.SetReference(branchRef)
}

// ListBranches 列出所有分支
func (s *Service) ListBranches(owner, name string) ([]string, error) {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}

	refs, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	var branches []string
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branches = append(branches, ref.Name().Short())
		}
		return nil
	})
	return branches, err
}

// ListCommits 列出提交历史
func (s *Service) ListCommits(owner, name, ref string) ([]*object.Commit, error) {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}

	refObj, err := repo.Reference(plumbing.ReferenceName(ref), true)
	if err != nil {
		return nil, err
	}

	commitsIter, err := repo.Log(&git.LogOptions{From: refObj.Hash()})
	if err != nil {
		return nil, err
	}
	defer commitsIter.Close()

	var commits []*object.Commit
	err = commitsIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	return commits, err
}

// CreateTag 创建标签
func (s *Service) CreateTag(owner, name, tagName, ref, message string) error {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return err
	}

	// 获取引用对应的提交
	refObj, err := repo.Reference(plumbing.ReferenceName(ref), true)
	if err != nil {
		return err
	}

	// 直接创建标签引用
	tagRef := plumbing.NewHashReference(
		plumbing.ReferenceName("refs/tags/"+tagName),
		refObj.Hash(),
	)

	return repo.Storer.SetReference(tagRef)
}

// DeleteTag 删除标签
func (s *Service) DeleteTag(owner, name, tagName string) error {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return err
	}

	tagRef := plumbing.ReferenceName("refs/tags/" + tagName)
	return repo.Storer.RemoveReference(tagRef)
}

// ListTags 列出所有标签
func (s *Service) ListTags(owner, name string) ([]string, error) {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}

	refs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var tags []string
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, ref.Name().Short())
		return nil
	})
	return tags, err
}

// GetTag 获取标签信息
func (s *Service) GetTag(owner, name, tagName string) (*object.Tag, error) {
	repo, err := s.GetRepo(owner, name)
	if err != nil {
		return nil, err
	}

	tagRef := plumbing.ReferenceName("refs/tags/" + tagName)
	_, err = repo.Reference(tagRef, true)
	if err != nil {
		return nil, err
	}

	// 对于轻量级标签，返回 nil
	return nil, nil
}

// ForkRepo 复制一个仓库
func (s *Service) ForkRepo(ctx context.Context, sourceOwner, sourceName, targetOwner, targetName string) error {
	// 获取源仓库路径
	sourcePath := s.getRepoPath(sourceOwner, sourceName)
	// 获取目标仓库路径
	targetPath := s.getRepoPath(targetOwner, targetName)

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("创建目标仓库目录失败: %w", err)
	}

	// 克隆仓库
	_, err := git.PlainClone(targetPath, true, &git.CloneOptions{
		URL:               sourcePath,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     plumbing.HEAD,
	})
	if err != nil {
		return fmt.Errorf("克隆仓库失败: %w", err)
	}

	return nil
}

// LFSStorePath 获取 LFS 存储路径
func (s *Service) LFSStorePath() string {
	return filepath.Join(s.repoBasePath, "lfs")
}

// EnsureLFSStore 确保 LFS 存储目录存在
func (s *Service) EnsureLFSStore() error {
	lfsPath := s.LFSStorePath()
	return os.MkdirAll(lfsPath, 0755)
}

// GetLFSObjectPath 获取 LFS 对象存储路径
func (s *Service) GetLFSObjectPath(oid string) string {
	// LFS 对象存储格式: lfs/xx/xx/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	if len(oid) < 4 {
		return ""
	}
	prefix1 := oid[:2]
	prefix2 := oid[2:4]
	return filepath.Join(s.LFSStorePath(), prefix1, prefix2, oid)
}

// UploadLFSObject 上传 LFS 对象
func (s *Service) UploadLFSObject(oid string, size int64, content []byte) error {
	if err := s.EnsureLFSStore(); err != nil {
		return err
	}

	objectPath := s.GetLFSObjectPath(oid)
	if objectPath == "" {
		return fmt.Errorf("invalid oid")
	}

	// 确保对象目录存在
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return fmt.Errorf("创建 LFS 对象目录失败: %w", err)
	}

	// 写入对象
	if err := os.WriteFile(objectPath, content, 0644); err != nil {
		return fmt.Errorf("写入 LFS 对象失败: %w", err)
	}

	return nil
}

// DownloadLFSObject 下载 LFS 对象
func (s *Service) DownloadLFSObject(oid string) ([]byte, error) {
	objectPath := s.GetLFSObjectPath(oid)
	if objectPath == "" {
		return nil, fmt.Errorf("invalid oid")
	}

	content, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, fmt.Errorf("读取 LFS 对象失败: %w", err)
	}

	return content, nil
}
