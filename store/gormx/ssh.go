package gormx

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

type SSHConfig struct {
	SshHost string `json:"sshHost" toml:"sshHost"`
	SshUser string `json:"sshUser" toml:"sshUser"`
	SshPk   string `json:"sshPk" toml:"sshPK"`
	SshPort int    `json:"sshPort" toml:"sshPort"`
}

type SSHDialer struct {
	client *ssh.Client
}

func (d *SSHDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return d.client.Dial("tcp", addr)
}

func (d *SSHDialer) Register() {
	mysql.RegisterDialContext("ssh", d.Dial)
}

func CreateSSHConn(cnf *SSHConfig) (*ssh.Client, error) {
	var pkBuf, err = ioutil.ReadFile(cnf.SshPk)
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer
	signer, err = ssh.ParsePrivateKey(pkBuf)
	if err != nil {
		return nil, err
	}

	var snf = &ssh.ClientConfig{
		User: cnf.SshUser,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	var sshCli *ssh.Client
	var sshUrl = fmt.Sprintf("%s:%d", cnf.SshHost, cnf.SshPort)
	sshCli, err = ssh.Dial("tcp", sshUrl, snf)
	if err != nil {
		return nil, err
	}
	return sshCli, nil
}
