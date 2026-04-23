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
	"time"
	"os/exec"

	"laima/internal/user/app"
	"golang.org/x/crypto/ssh"
)

// SSHServer SSH 服务器
type SSHServer struct {
	Addr         string
	HostKeyPath  string
	RepoBasePath string
	gitService   *Service
	userService  app.UserService
	hostKey      ssh.Signer // 缓存的主机密钥
}

// NewSSHServer 创建 SSH 服务器实例
func NewSSHServer(addr, hostKeyPath, repoBasePath string, gitService *Service, userService app.UserService) *SSHServer {
	return &SSHServer{
		Addr:         addr,
		HostKeyPath:  hostKeyPath,
		RepoBasePath: repoBasePath,
		gitService:   gitService,
		userService:  userService,
	}
}

// Start 启动 SSH 服务器
func (s *SSHServer) Start(ctx context.Context) error {
	log.Printf("正在启动 SSH 服务器...")

	// 加载主机密钥
	hostKey, err := s.loadHostKey()
	if err != nil {
		log.Printf("错误: 加载主机密钥失败: %v", err)
		return fmt.Errorf("加载主机密钥失败: %w", err)
	}
	log.Printf("主机密钥加载成功")

	// 配置 SSH 服务器
	config := &ssh.ServerConfig{
		// 添加公钥认证
		PublicKeyCallback: s.handlePublicKeyAuth,
		// 禁用密码认证
		PasswordCallback: nil,
	}
	config.AddHostKey(hostKey)
	log.Printf("SSH 服务器配置完成")

	// 监听端口
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Printf("错误: 监听端口 %s 失败: %v", s.Addr, err)
		return fmt.Errorf("监听失败: %w", err)
	}
	log.Printf("成功监听端口 %s", s.Addr)

	// 启动服务器
	go s.serve(listener, config, ctx)

	log.Printf("SSH 服务器启动在 %s", s.Addr)
	return nil
}

// loadHostKey 加载主机密钥
func (s *SSHServer) loadHostKey() (ssh.Signer, error) {
	// 如果主机密钥已缓存，直接返回
	if s.hostKey != nil {
		return s.hostKey, nil
	}

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

	// 缓存主机密钥
	s.hostKey = signer

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

	// 清除主机密钥缓存，确保下次加载时使用新的密钥
	s.hostKey = nil

	return nil
}

// handlePublicKeyAuth 处理 SSH 公钥认证
func (s *SSHServer) handlePublicKeyAuth(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
	// 计算公钥指纹
	fingerprint := ssh.FingerprintSHA256(pubKey)
	log.Printf("SSH认证尝试: 用户 %s, 指纹 %s", c.User(), fingerprint)

	// 如果没有用户服务，暂时允许所有认证
	if s.userService == nil {
		log.Printf("警告: 没有用户服务，允许所有 SSH 认证")
		return &ssh.Permissions{
			Extensions: map[string]string{
				"user": c.User(),
			},
		}, nil
	}

	// 尝试根据指纹查找 SSH 密钥
	sshKey, err := s.userService.GetSSHKeyByFingerprint(fingerprint)
	if err != nil {
		log.Printf("SSH认证失败: %v", err)
		return nil, fmt.Errorf("invalid public key")
	}

	// 验证公钥是否匹配
	savedKey, err := ssh.ParsePublicKey([]byte(sshKey.Key))
	if err != nil {
		log.Printf("解析保存的 SSH 密钥失败: %v", err)
		return nil, fmt.Errorf("invalid saved public key")
	}

	// 比较公钥
	if string(pubKey.Marshal()) != string(savedKey.Marshal()) {
		log.Printf("SSH 密钥不匹配")
		return nil, fmt.Errorf("public key mismatch")
	}

	log.Printf("SSH认证成功: 用户 ID %d", sshKey.UserID)

	// 返回认证成功的权限
	return &ssh.Permissions{
		Extensions: map[string]string{
			"user_id": fmt.Sprintf("%d", sshKey.UserID),
			"user":    c.User(),
		},
	}, nil
}

// serve 处理 SSH 连接
func (s *SSHServer) serve(listener net.Listener, config *ssh.ServerConfig, ctx context.Context) {
	defer listener.Close()

	// 创建工作池
	workerPool := make(chan struct{}, 100) // 限制并发连接数为100

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

			// 限制并发连接数
			workerPool <- struct{}{}

			// 处理连接
			go func() {
				defer func() {
					<-workerPool
				}()
				s.handleConnection(conn, config)
			}()
		}
	}
}

// handleConnection 处理单个 SSH 连接
func (s *SSHServer) handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	// 记录连接信息
	clientAddr := conn.RemoteAddr().String()
	log.Printf("收到新的 SSH 连接来自 %s", clientAddr)

	// 进行 SSH 握手
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("错误: SSH 握手失败来自 %s: %v", clientAddr, err)
		return
	}
	defer sshConn.Close()

	// 记录认证成功的用户信息
	log.Printf("SSH 连接认证成功来自 %s，用户: %s", clientAddr, sshConn.User())

	// 忽略全局请求
	go ssh.DiscardRequests(reqs)

	// 处理通道
	for newChan := range chans {
		log.Printf("收到通道请求类型: %s 来自 %s", newChan.ChannelType(), clientAddr)
		if newChan.ChannelType() != "session" {
			log.Printf("拒绝非 session 通道类型: %s 来自 %s", newChan.ChannelType(), clientAddr)
			newChan.Reject(ssh.UnknownChannelType, "只支持 session 通道")
			continue
		}

		// 接受通道
		ch, reqs, err := newChan.Accept()
		if err != nil {
			log.Printf("错误: 接受通道失败来自 %s: %v", clientAddr, err)
			continue
		}

		// 处理会话
		go func() {
			defer ch.Close()
			s.handleSession(ch, reqs)
			log.Printf("SSH 会话结束来自 %s", clientAddr)
		}()
	}
	log.Printf("SSH 连接关闭来自 %s", clientAddr)
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
			} else {
				log.Printf("警告: 拒绝执行非 Git 命令: %s", cmd)
				io.WriteString(ch, "错误: 只支持 Git 命令\n")
			}

			// 回复请求
			if err := req.Reply(true, nil); err != nil {
				log.Printf("错误: 回复请求失败: %v", err)
			}
		} else {
			log.Printf("忽略请求类型: %s", req.Type)
			if err := req.Reply(false, nil); err != nil {
				log.Printf("错误: 回复请求失败: %v", err)
			}
		}
	}
}

// handleGitCommand 处理 Git 命令
func (s *SSHServer) handleGitCommand(ch ssh.Channel, cmd string) {
	// 解析 Git 命令
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		log.Printf("错误: 无效的 Git 命令: %s", cmd)
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
		log.Printf("错误: 无效的仓库路径: %s", repoPath)
		io.WriteString(ch, "错误: 无效的仓库路径\n")
		return
	}

	owner := parts[len(parts)-2]
	repoName := parts[len(parts)-1]

	// 获取仓库路径
	fullRepoPath := s.gitService.getRepoPath(owner, repoName)
	log.Printf("处理 Git 命令 %s 仓库 %s/%s，路径: %s", gitCmd, owner, repoName, fullRepoPath)

	// 检查仓库是否存在
	if _, err := os.Stat(fullRepoPath); os.IsNotExist(err) {
		log.Printf("错误: 仓库不存在: %s/%s", owner, repoName)
		io.WriteString(ch, "错误: 仓库不存在\n")
		return
	} else if err != nil {
		log.Printf("错误: 检查仓库存在性失败: %v", err)
		io.WriteString(ch, fmt.Sprintf("错误: 检查仓库失败: %v\n", err))
		return
	}

	// 处理 Git 命令
	switch gitCmd {
	case "git-receive-pack":
		log.Printf("执行 git-receive-pack 命令 for %s/%s", owner, repoName)
		s.handleGitReceivePack(ch, fullRepoPath)
		log.Printf("git-receive-pack 命令执行完成 for %s/%s", owner, repoName)
	case "git-upload-pack":
		log.Printf("执行 git-upload-pack 命令 for %s/%s", owner, repoName)
		s.handleGitUploadPack(ch, fullRepoPath)
		log.Printf("git-upload-pack 命令执行完成 for %s/%s", owner, repoName)
	default:
		log.Printf("错误: 不支持的 Git 命令: %s", gitCmd)
		io.WriteString(ch, fmt.Sprintf("错误: 不支持的 Git 命令: %s\n", gitCmd))
	}
}

// handleGitReceivePack 处理 git-receive-pack 命令（推送）
func (s *SSHServer) handleGitReceivePack(ch ssh.Channel, repoPath string) {
	// 设置 10 分钟超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 使用系统 Git 命令处理
	cmd := exec.CommandContext(ctx, "git", "receive-pack", repoPath)
	cmd.Stdin = ch
	cmd.Stdout = ch
	cmd.Stderr = ch

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			io.WriteString(ch, "错误: 命令执行超时\n")
		} else {
			io.WriteString(ch, fmt.Sprintf("错误: %v\n", err))
		}
	}
}

// handleGitUploadPack 处理 git-upload-pack 命令（拉取）
func (s *SSHServer) handleGitUploadPack(ch ssh.Channel, repoPath string) {
	// 设置 10 分钟超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 使用系统 Git 命令处理
	cmd := exec.CommandContext(ctx, "git", "upload-pack", repoPath)
	cmd.Stdin = ch
	cmd.Stdout = ch
	cmd.Stderr = ch

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			io.WriteString(ch, "错误: 命令执行超时\n")
		} else {
			io.WriteString(ch, fmt.Sprintf("错误: %v\n", err))
		}
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
