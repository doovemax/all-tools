package cmd

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh/knownhosts"

	"golang.org/x/crypto/ssh"

	"github.com/doovemax/all-tool/conf"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(copyKey)
	KeyFile = copyKey.PersistentFlags().StringP("key", "k", "$HOME/.ssh/id_rsa.pub", "指定key文件")
}

var copyKey = &cobra.Command{
	Use:                        "copykey",
	Aliases:                    nil,
	SuggestFor:                 nil,
	Short:                      "",
	Long:                       "",
	Example:                    "",
	ValidArgs:                  nil,
	Args:                       nil,
	ArgAliases:                 nil,
	BashCompletionFunction:     "",
	Deprecated:                 "",
	Hidden:                     false,
	Annotations:                nil,
	Version:                    "",
	PersistentPreRun:           nil,
	PersistentPreRunE:          nil,
	PreRun:                     nil,
	PreRunE:                    nil,
	Run:                        readKey,
	RunE:                       nil,
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
	TraverseChildren:           true,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
}

func readKey(c *cobra.Command, args []string) {
	var (
		idRsaPub []byte
		keyfile  string
		err      error
	)
	conf.Hosts, err = conf.ReadConf(*ConfFile)
	if err != nil {
		logrus.Errorf("配置文件读取错误：%v", err)
	}

	if c.Flag("key").Changed {
		keyfile = *KeyFile
	} else {
		keyfile = strings.Replace(*KeyFile, "$HOME", os.Getenv("HOME"), -1)
	}
	idRsaPub, err = ioutil.ReadFile(keyfile)
	if err != nil {
		logrus.Fatalf("ssh public key error: %v", err)
		logrus.Exit(2)
	}
	logrus.Infoln("--------------go go go------------------")
	for index, host := range conf.Hosts {

		var LimitHostsList []string
		var re *regexp.Regexp
		var flag = 0
		if *LimitHosts != "" && *RegexFlag == false {
			LimitHostsList = strings.Split(*LimitHosts, ",")
			flag = 1
		} else if *LimitHosts != "" && *RegexFlag == true {
			re, err = regexp.Compile(*LimitHosts)
			if err != nil {
				return
			}
			flag = 2
		}

		if flag == 1 {
			for _, h := range LimitHostsList {
				if h == host.Name {
					goto mach
				}

			}
			continue

		} else if flag == 2 {
			if re.MatchString(host.Name) {
				goto mach
			}
			continue
		}

	mach:

		session, err := host.SshClient()
		if err != nil {
			logrus.Errorf("ssh session error:  %v", err)
			continue
		}
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		err = session.Run(fmt.Sprintf("echo '%s' >> $HOME/.ssh/authorized_keys && chmod 600 $HOME/.ssh/authorized_keys ", string(idRsaPub)))
		if err != nil {
			logrus.Errorf("line: %d %-16v%-6v%v: %s\n", index+1, host.IP, host.Port, host.User, err.Error())
		} else {

			logrus.Infof("line: %d %-16v%-6v%v: %s\n", index+1, host.IP, host.Port, host.User, "Successing")

			IP, err := getIP()
			if err != nil {
				logrus.Errorf("无法获取 IP 地址: %v\n", err)
			} else {
				line, err := getKnownHost(idRsaPub, IP)

				if err != nil {
					logrus.Errorf("生成 Knownhosts 失败: %v\n", err)
				} else {
					err = WriteKnownHost(line)
					if err != nil {
						logrus.Errorf("Known_hosts 写入失败：%v", err)
					}
				}
			}

		}

	}
	logrus.Infoln("--------------over------------------")
}

func getKnownHost(pubKey []byte, IP string) (line string, err error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(pubKey)
	if err != nil {
		return "", err
	}
	ip := strings.Split(IP, ".")

	knownHostsLine := knownhosts.Line(ip, pub)
	return knownHostsLine, err

}

func getIP() (IP string, err error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), nil
					}
				}
			}

		}
	}
	return
}

func WriteKnownHost(line string) (err error) {
	var f *os.File
	home := os.Getenv("HOME")
	_, err = os.Stat(home + "/.ssh/known_hosts")
	if err == nil {
		f, err = os.OpenFile(home+"/.ssh/known_hosts", os.O_WRONLY|os.O_APPEND, 644)
		if err != nil {
			return err
		}
	} else {
		f, err = os.OpenFile(home+"/.ssh/known_hosts", os.O_CREATE|os.O_WRONLY, 644)
		if err != nil {
			return err
		}
	}

	_, err = f.WriteString(line + "\n")
	if err != nil {
		return
	}

	return nil
}
