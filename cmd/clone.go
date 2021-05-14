package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const gitBase = "https://codeload.github.com"
const tempFileName = "roverclone.tar.gz"

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Fetch supporting artifacts such as landingzones from GitHub",
	Long: `Pull down repos from GitHub and extracts them in well defined way.
Git is not required`,

	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		strip, _ := cmd.Flags().GetInt("strip")
		dest, _ := cmd.Flags().GetString("dest")
		folder, _ := cmd.Flags().GetString("folder")
		runClone(repo, branch, strip, dest, folder)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringP("repo", "r", "azure/caf-terraform-landingzones", "Which repo on GitHub to clone")
	cloneCmd.Flags().StringP("branch", "b", "master", "Which branch to clone")
	cloneCmd.Flags().IntP("strip", "s", 1, "Levels to strip from repo hierarchy, best left as 1")
	cloneCmd.Flags().StringP("dest", "d", "./landingzones", "Where to place output")
	cloneCmd.Flags().StringP("folder", "f", "", "Extract a sub-folder from the repo")
}

func runClone(repo string, branch string, strip int, dest string, subFolder string) {
	command.EnsureDirectory(viper.GetString("tempDir"))
	command.RemoveDirectory(dest)
	command.EnsureDirectory(dest)
	color.Green("Running clone operation. Will fetch %s branch of %s and place into %s", branch, repo, dest)

	tempDir := viper.GetString("tempDir")
	cloneUrl := fmt.Sprintf("%s/%s/tar.gz/%s", gitBase, repo, branch)
	tarFile := fmt.Sprintf("%s/%s", tempDir, tempFileName)

	projParts := strings.Split(repo, "/")
	if len(projParts) < 2 {
		cobra.CheckErr("repo '" + repo + "' is invalid, should be in form `org/name`")
	}
	projName := projParts[1]

	folderName := fmt.Sprintf("%s-%s/%s", projName, branch, subFolder)

	cmd := command.NewCommand("curl", []string{
		cloneUrl, "--fail", "--silent", "--show-error", "-o", tarFile,
	})
	err := cmd.Execute(true)
	cobra.CheckErr(err)

	cmd = command.NewCommand("tar", []string{
		"-zxvf", tarFile, "--strip-components", strconv.Itoa(strip), "-C", dest, folderName,
	})
	err = cmd.Execute(true)
	cobra.CheckErr(err)
}
