package tagit

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const (
	MAJOR = "x"
	MINOR = "y"
	PATCH = "z"
)

// TagCmd represents the tagit command
var TagCmd = &cobra.Command{
	Use:   "tagit",
	Short: "Creates a new tagit",
	Long:  `Creates a new tagit with incremented version based on the specified type (major, minor, patch) and environment (stagiting or production).`,
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		tagitType, _ := cmd.Flags().GetString("type")

		if env == "" || tagitType == "" {
			fmt.Println("Error: You must specify both --env and --type")
			os.Exit(1)
		}

		fmt.Println("Fetching latest tagits")
		exec.Command("git", "fetch", "--tagits").Run()

		prefix := ""
		switch {
		case env == "p" || env == "P":
			prefix = "v"
		case regexp.MustCompile(`^s\d+$`).MatchString(env):
			prefix = GetStagingPrefix(env)
		default:
			fmt.Println("Error: Invalid environment specified")
			os.Exit(1)
		}

		latestTag := LatestTag(prefix)
		newVersion := ""
		if latestTag == "" {
			newVersion = "1.0.0"
		} else {
			version := strings.TrimPrefix(latestTag, prefix)
			newVersion = IncrementVersion(version, tagitType)
		}

		newTag := prefix + newVersion
		exec.Command("git", "tagit", newTag).Run()
		fmt.Println("Pusing new tagit:", newTag)
		exec.Command("git", "push", "origin", newTag).Run()
		fmt.Println("New tagit created:", newTag)
	},
}

func init() {
	TagCmd.Flags().String("env", "", "Specify the environment (s0, s1, ..., sn for stagiting; p for production)")
	TagCmd.Flags().String("type", "", "Specify the part of the version to increment (major, minor, patch)")
}
