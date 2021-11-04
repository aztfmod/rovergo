//
// Rover - Landing zone fetch
// * This fetches the official CAF landing zones from GitHub, replaces the old --clone feature
// * Ben C, May 2021
//

package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/rover"
	"github.com/spf13/cobra"
)

const gitBase = "https://codeload.github.com"
const tempFileName = "roverfetch.tar.gz"

var lzFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch supporting artifacts such as landingzones from GitHub",
	Long: `Pull down landingzone repos from GitHub and extracts them in well defined way.
Git is not required`,

	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		strip, _ := cmd.Flags().GetInt("strip")
		dest, _ := cmd.Flags().GetString("dest")
		folder, _ := cmd.Flags().GetString("folder")
		runFetch(repo, branch, strip, dest, folder)
	},
}

func init() {
	landingzoneCmd.AddCommand(lzFetchCmd)
	lzFetchCmd.Flags().StringP("repo", "r", "azure/caf-terraform-landingzones", "Which repo on GitHub to fetch")
	lzFetchCmd.Flags().StringP("branch", "b", "master", "Which branch to fetch")
	lzFetchCmd.Flags().IntP("strip", "s", 1, "Levels to strip from repo hierarchy, best left as 1")
	lzFetchCmd.Flags().StringP("dest", "d", "./landingzones", "Where to place output")
	lzFetchCmd.Flags().StringP("folder", "f", "", "Extract a sub-folder from the repo")
}

func runFetch(repo string, branch string, strip int, dest string, subFolder string) {

	homeDir, homeErr := rover.HomeDirectory()
	cobra.CheckErr(homeErr)
	tempDir, dirErr := ioutil.TempDir(homeDir, "fetchops*")
	cobra.CheckErr(dirErr)
	command.EnsureDirectory(tempDir)
	command.RemoveDirectory(dest)
	command.EnsureDirectory(dest)
	console.Infof("Running fetch operation. Will download %s branch of %s and place into %s\n", branch, repo, dest)

	cloneURL := fmt.Sprintf("%s/%s/tar.gz/%s", gitBase, repo, branch)
	tarFile := path.Join(tempDir, tempFileName)

	projParts := strings.Split(repo, "/")
	if len(projParts) < 2 {
		cobra.CheckErr("repo '" + repo + "' is invalid, should be in form `org/name`")
	}
	projName := projParts[1]

	folderName := fmt.Sprintf("%s-%s/%s", projName, branch, subFolder)

	cmd := command.NewCommand("curl", []string{
		cloneURL, "--fail", "--silent", "--show-error", "-o", tarFile,
	})
	err := cmd.Execute()
	cobra.CheckErr(err)

	cmd = command.NewCommand("tar", []string{
		"-zxvf", tarFile, "--strip-components", strconv.Itoa(strip), "-C", dest, folderName,
	})
	err = cmd.Execute()
	cobra.CheckErr(err)
}
