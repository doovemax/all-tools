package conf

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	Hosts []*hostInfo
)

// 文件格式:  name  user  ip port password
type hostInfo struct {
	Name     string
	User     string
	IP       string
	Port     string
	Password string
}

func (this *hostInfo) SshClient() (session *ssh.Session, err error) {

	client, err := ssh.Dial("tcp", this.IP+":"+this.Port, &ssh.ClientConfig{
		Config: ssh.Config{},
		User:   this.User,
		Auth:   []ssh.AuthMethod{ssh.Password(this.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		BannerCallback:    nil,
		ClientVersion:     "",
		HostKeyAlgorithms: nil,
		Timeout:           0,
	})
	// defer client.Close()
	if err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return
	}
	// defer session.Close()
	return
}

func ReadConf(filepath string) (Hosts []*hostInfo, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	fb, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	hostList := bytes.Split(fb, []byte("\n"))

	for index, host := range hostList {
		info := bytes.Fields(host)
		if len(info) != 5 {

			return nil, errors.New(fmt.Sprintf("格式不对：Line %v %v", index+1, string(host)))
		}

		// 过滤 hosts

		Hosts = append(Hosts, &hostInfo{
			Name:     string(info[0]),
			User:     string(info[1]),
			IP:       string(info[2]),
			Port:     string(info[3]),
			Password: string(info[4]),
		})

	}
	return
}
