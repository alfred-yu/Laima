package git

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

// SSHServer SSH 服务器
type SSHServer struct {
	Addr         string
	HostKeyPath  string
	RepoBasePath string
	gitService   *Service
}

// NewSSHServer 创建 SSH 服务器实例
func NewSSHServer(addr, hostKeyPath, repoBasePath string, gitService *Service) *SSHServer {
	return &SSHServer{
		Addr:         addr,
		HostKeyPath:  hostKeyPath,
		RepoBasePath: repoBasePath,
		gitService:   gitService,
	}
}

// Start 启动 SSH 服务器
func (s *SSHServer) Start(ctx context.Context) error {
	// 加载主机密钥
	hostKey, err := s.loadHostKey()
	if err != nil {
		return fmt.Errorf("加载主机密钥失败: %w", err)
	}

	// 配置 SSH 服务器
	config := &ssh.ServerConfig{
		NoClientAuth: true, // 暂时不进行客户端认证，后续可以添加
	}
	config.AddHostKey(hostKey)

	// 监听端口
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("监听失败: %w", err)
	}

	// 启动服务器
	go s.serve(listener, config, ctx)

	log.Printf("SSH 服务器启动在 %s", s.Addr)
	return nil
}

// loadHostKey 加载主机密钥
func (s *SSHServer) loadHostKey() (ssh.Signer, error) {
	// 如果主机密钥文件不存在，生成一个新的
	if _, err := os.Stat(s.HostKeyPath); os.IsNotExist(err) {
		if err := s.generateHostKey(); err != nil {
			return nil, err
		}
	}

	// 读取主机密钥
	hostKeyBytes, err := os.ReadFile(s.HostKeyPath)
	if err != nil {
		return nil, err
	}

	// 解析主机密钥
	signer, err := ssh.ParsePrivateKey(hostKeyBytes)
	if err != nil {
		return nil, err
	}

	return signer, nil
}

// generateHostKey 生成主机密钥
func (s *SSHServer) generateHostKey() error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(s.HostKeyPath), 0755); err != nil {
		return err
	}

	// 生成新的 RSA 密钥对
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 保存私钥
	hostKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err := os.WriteFile(s.HostKeyPath, hostKeyBytes, 0600); err != nil {
		return err
	}

	return nil
}

// serve 处理 SSH 连接
func (s *SSHServer) serve(listener net.Listener, config *ssh.ServerConfig, ctx context.Context) {
	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 接受新连接
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("接受连接失败: %v", err)
				continue
			}

			// 处理连接
			go s.handleConnection(conn, config)
		}
	}
}

// handleConnection 处理单个 SSH 连接
func (s *SSHServer) handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	// 进行 SSH 握手
	_, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("SSH 握手失败: %v", err)
		return
	}

	// 忽略全局请求
	go ssh.DiscardRequests(reqs)

	// 处理通道
	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			newChan.Reject(ssh.UnknownChannelType, "只支持 session 通道")
			continue
		}

		// 接受通道
		ch, reqs, err := newChan.Accept()
		if err != nil {
			log.Printf("接受通道失败: %v", err)
			continue
		}

		// 处理会话
		go s.handleSession(ch, reqs)
	}
}

// handleSession 处理 SSH 会话
func (s *SSHServer) handleSession(ch ssh.Channel, reqs <-chan *ssh.Request) {
	defer ch.Close()

	for req := range reqs {
		if req.Type == "exec" {
			// 解析命令
			cmd := string(req.Payload[4:])
			log.Printf("执行命令: %s", cmd)

			// 处理 Git 命令
			if strings.HasPrefix(cmd, "git-") {
				s.handleGitCommand(ch, cmd)
			}

			req.Reply(true, nil)
		}
	}
}

// handleGitCommand 处理 Git 命令
func (s *SSHServer) handleGitCommand(ch ssh.Channel, cmd string) {
	// 解析 Git 命令
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		io.WriteString(ch, "错误: 无效的 Git 命令\n")
		return
	}

	gitCmd := parts[0]
	repoPath := parts[1]

	// 清理仓库路径
	repoPath = strings.TrimPrefix(repoPath, "'")
	repoPath = strings.TrimSuffix(repoPath, "'")
	repoPath = strings.TrimSuffix(repoPath, ".git")

	// 解析所有者和仓库名
	parts = strings.Split(repoPath, "/")
	if len(parts) < 2 {
		io.WriteString(ch, "错误: 无效的仓库路径\n")
		return
	}

	owner := parts[len(parts)-2]
	repoName := parts[len(parts)-1]

	// 获取仓库路径
	fullRepoPath := s.gitService.getRepoPath(owner, repoName)

	// 检查仓库是否存在
	if _, err := os.Stat(fullRepoPath); os.IsNotExist(err) {
		io.WriteString(ch, "错误: 仓库不存在\n")
		return
	}

	// 处理 Git 命令
	switch gitCmd {
	case "git-receive-pack":
		s.handleGitReceivePack(ch, fullRepoPath)
	case "git-upload-pack":
		s.handleGitUploadPack(ch, fullRepoPath)
	default:
		io.WriteString(ch, fmt.Sprintf("错误: 不支持的 Git 命令: %s\n", gitCmd))
	}
}

// handleGitReceivePack 处理 git-receive-pack 命令（推送）
func (s *SSHServer) handleGitReceivePack(ch ssh.Channel, repoPath string) {
	// 使用系统 Git 命令处理
	cmd := exec.Command("git", "receive-pack", repoPath)
	cmd.Stdin = ch
	cmd.Stdout = ch
	cmd.Stderr = ch

	if err := cmd.Run(); err != nil {
		io.WriteString(ch, fmt.Sprintf("错误: %v\n", err))
	}
}

// handleGitUploadPack 处理 git-upload-pack 命令（拉取）
func (s *SSHServer) handleGitUploadPack(ch ssh.Channel, repoPath string) {
	// 使用系统 Git 命令处理
	cmd := exec.Command("git", "upload-pack", repoPath)
	cmd.Stdin = ch
	cmd.Stdout = ch
	cmd.Stderr = ch

	if err := cmd.Run(); err != nil {
		io.WriteString(ch, fmt.Sprintf("错误: %v\n", err))
	}
}

// GitTransportService Git 传输服务
type GitTransportService struct {
	repoPath string
}

// NewGitTransportService 创建 Git 传输服务
func NewGitTransportService(repoPath string) *GitTransportService {
	return &GitTransportService{
		repoPath: repoPath,
	}
}

// NewReceivePackSession 创建接收包会话
func (s *GitTransportService) NewReceivePackSession(ctx context.Context, endpoint interface{}, auth interface{}) (interface{}, error) {
	return nil, nil
}

// NewUploadPackSession 创建上传包会话
func (s *GitTransportService) NewUploadPackSession(ctx context.Context, endpoint interface{}, auth interface{}) (interface{}, error) {
	return nil, nil
}
