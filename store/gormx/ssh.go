package gormx

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host string `json:"host" toml:"host"`
	User string `json:"user" toml:"user"`
	Pk   string `json:"pk" toml:"pk"`
	Port int    `json:"port" toml:"port"`
}

type SSHDialer struct {
	client *ssh.Client
}

func (d *SSHDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return d.client.DialContext(ctx, "tcp", addr)
}

func (d *SSHDialer) Register() {
	mysql.RegisterDialContext("ssh", d.Dial)
}

func CreateSSHConn(cnf *SSHConfig) (*ssh.Client, error) {
	var pkBuf, err = os.ReadFile(cnf.Pk)
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer
	signer, err = ssh.ParsePrivateKey(pkBuf)
	if err != nil {
		return nil, err
	}

	var snf = &ssh.ClientConfig{
		User: cnf.User,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(_ string, _ net.Addr, _ ssh.PublicKey) error {
			return nil
		},
	}
	var sshCli *ssh.Client
	var sshURL = fmt.Sprintf("%s:%d", cnf.Host, cnf.Port)
	sshCli, err = ssh.Dial("tcp", sshURL, snf)
	if err != nil {
		return nil, err
	}
	return sshCli, nil
}
