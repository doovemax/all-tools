package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
	for index, host := range conf.Hosts {
		session, err := host.SshClient()
		if err != nil {
			logrus.Errorf("ssh session error:  %v", err)
			continue
		}
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		err = session.Run(fmt.Sprintf("echo '%s' >> $HOME/.ssh/authorized_keys && chmod 600 $HOME/.ssh/authorized_keys ", string(idRsaPub)))
		if err != nil {
			logrus.Printf("%d) %-16v%-6v%v: %s\n", index+1, host.IP, host.Port, host.User, err.Error())
		} else {
			logrus.Printf("%d) %-16v%-6v%v: %s\n", index+1, host.IP, host.Port, host.User, "Successing")
		}

	}
}
