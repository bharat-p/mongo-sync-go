package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	VersionString string
	Commit        string
	BuildDate     string
	versionString string
	revString     string
	buildDate     string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `All software has versions. This is mongo-sync-go's.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printVersion()
		return nil
	},
}

func printVersion() {
	if len(Commit) > 7 {
		Commit = Commit[0:7]
	}
	fmt.Printf("%-15s %s\n", "Version:", VersionString)
	fmt.Printf("%-15s %s\n", "Go version:", runtime.Version())
	fmt.Printf("%-15s %s\n", "Git commit:", Commit)
	fmt.Printf("%-15s %s\n", "Built:", BuildDate)
	fmt.Printf("%-15s %s/%s\n", "OS/Arch:", runtime.GOOS, runtime.GOARCH)

}

func init() {
	RootCmd.AddCommand(versionCmd)
}
