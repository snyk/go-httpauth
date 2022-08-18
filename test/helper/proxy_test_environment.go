package helper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type ProxyTestEmvironment struct {
	repositoryBaseDirectory string
}

func NewProxyTestEnvironment(RepositoryBaseDirectory string) *ProxyTestEmvironment {
	return &ProxyTestEmvironment{repositoryBaseDirectory: RepositoryBaseDirectory}
}

func (p *ProxyTestEmvironment) FixturePath() string {
	path, _ := filepath.Abs(filepath.Join(p.repositoryBaseDirectory, "test", "fixtures"))
	path, _ = filepath.Abs(path)
	return path
}

func (p *ProxyTestEmvironment) DockerComposeFile() string {
	path := filepath.Join(p.FixturePath(), "squid_environment", "docker-compose.yml")
	path, _ = filepath.Abs(path)
	return path
}

func (p *ProxyTestEmvironment) ScriptsPath() string {
	path := filepath.Join(p.FixturePath(), "squid_environment", "scripts")
	path, _ = filepath.Abs(path)
	return path
}

func (p *ProxyTestEmvironment) ConfigFile() string {
	path := filepath.Join(p.ScriptsPath(), "krb5.conf")
	return path
}

func (p *ProxyTestEmvironment) CacheFile() string {
	path := filepath.Join(p.ScriptsPath(), "krb5_cache")
	return path
}

func (p *ProxyTestEmvironment) HasDockerInstalled() bool {
	result := false
	cmd := exec.Command("docker-compose", "--version")
	err := cmd.Run()
	if err == nil {
		result = true
	}
	return result
}

func (p *ProxyTestEmvironment) RunTestUsingExternalAuthLibrary() bool {
	result := false
	os := runtime.GOOS
	if os == "windows" {
		result = true
	} else {
		result = p.HasDockerInstalled()
	}
	return result
}

func (p *ProxyTestEmvironment) RunDockerCompose(arg ...string) {
	cmd := exec.Command("docker-compose", arg...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "HTTP_PROXY_PORT=3128")
	cmd.Env = append(cmd.Env, "PROXY_HOSTNAME=proxy.snyk.local")
	cmd.Env = append(cmd.Env, "CONTAINER_NAME=spnego_test")
	cmd.Env = append(cmd.Env, "SCRIPTS_PATH="+p.ScriptsPath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	run := func() {
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	run()
}

func (p *ProxyTestEmvironment) StartProxyEnvironment() {
	p.StopProxyEnvironment()
	p.RunDockerCompose("--file", p.DockerComposeFile(), "up", "--build", "--detach")

	config := filepath.Join(p.FixturePath(), "squid_environment", "scripts", "krb5.conf")
	waitForFile(config, time.Second*30)
}

func (p *ProxyTestEmvironment) StopProxyEnvironment() {
	p.RunDockerCompose("--file", p.DockerComposeFile(), "down")
	os.Remove(p.CacheFile())
	os.Remove(p.ConfigFile())
}

func waitForFile(filename string, timeout time.Duration) {
	start := time.Now()
	for {
		_, err := os.Stat(filename)
		if !os.IsNotExist(err) {
			break
		}

		if time.Since(start) >= timeout {
			fmt.Println("waitForFile() - timeout", filename)
			break
		}

		time.Sleep(time.Second)

	}
}
