package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/toolkits/file"
)

var (
	ConfFile *string
	KeyFile  *string
)

var RootCmd = &cobra.Command{
	Use:                        "",
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
	Run:                        nil,
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
	TraverseChildren:           false,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
}

func init() {
	str := file.Basename(os.Args[0])
	RootCmd.Use = str
	ConfFile = RootCmd.PersistentFlags().StringP("conf", "f", "./hosts.conf", "按行读取配置文件:name user ip port password\n")
}

//
// func globle(c *cobra.Command, args []string) {
// 	// fmt.Println("test")
// 	return
// }
