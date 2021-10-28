package landingzone

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aztfmod/rover/pkg/azure"
	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/spf13/cobra"

	"github.com/jstemmer/go-junit-report/formatter"
	"github.com/jstemmer/go-junit-report/parser"
)

type TestAction struct {
	ActionBase
}

func NewTestAction() *TestAction {
	return &TestAction{
		ActionBase{
			Name:        "Test",
			Type:        BuiltinCommand,
			Description: "Perform a Go Test",
		},
	}
}

func (ta *TestAction) Execute(o *Options) error {
	console.Info("Carrying out Rover Test")

	// Locate state storage account id
	if o.StateSubscription == "" {
		console.Infof("No State sub provided, trying to locate default sub")
		sub, err := azure.GetSubscription()

		if err != nil {
			return errors.New("can't locate state sub")
		}
		o.StateSubscription = sub.ID
	}

	storageID, err := azure.FindStorageAccount(o.Level, o.CafEnvironment, o.StateSubscription)
	if err != nil {

		console.Errorf("No state storage account found for environment '%s' and level %s", o.CafEnvironment, o.Level)
		return errors.New("can't test a landing zone without a state file storage account")
	}

	console.Infof("Located state storage account %s\n", storageID)

	// download tfstate file
	stateFilePath := path.Join(o.DataDir, "terraform.tfstate")
	err = azure.DownloadFileFromBlob(storageID, o.Workspace, o.StateName+".tfstate", stateFilePath)
	cobra.CheckErr(err)

	// Execute go test
	console.Infof("Execute tests in %s\n", o.TestPath)
	console.StartSpinner()

	testRes, err := RunGoTests(o.TestPath, o)
	console.StopSpinner()

	// create junit test report.
	if err == nil {
		console.Infof("Test execution compeleted!")
		_ = CreateJunitReport(testRes, "testReport.xml")
	} else {
		console.Errorf("Test execution errored! %s\n", err)
	}

	// delete the downloaded statefile
	// Remove files
	o.cleanUp()
	_ = os.Remove(stateFilePath)

	return nil
}

func RunGoTests(testFilePath string, o *Options) (cmdRes string, err error) {
	// export env
	os.Setenv("STATE_FILE_PATH", o.DataDir)
	os.Setenv("ARM_SUBSCRIPTION_ID", o.StateSubscription)
	os.Setenv("ENVIRONMENT", o.CafEnvironment)

	args := []string{"test"}
	args = append(args, "-v")

	if o.Stack != "" {
		args = append(args, "--tags", fmt.Sprintf("%s,%s", o.Level, o.Stack))
	} else {
		args = append(args, "--tags", o.Level)
	}

	// Set the path to the test path, this is where the tests are
	currDir, _ := os.Getwd()
	err = os.Chdir(testFilePath)
	if err != nil {
		console.Errorf("Error switching to test directory %s %s\n", testFilePath, err)
	}
	cmd := command.NewCommand("go", args)
	cmd.Silent = false

	err = cmd.Execute()
	if err != nil {
		console.Infof("Executed tests output logs %s\n", cmd.StdOut)

		return cmd.StdErr, err
	}

	console.Infof("Executed tests output logs %s\n", cmd.StdOut)
	err = os.Chdir(currDir)
	if err != nil {
		console.Errorf("Error switching to back to working directory %s %s\n", currDir, err)
	}
	return cmd.StdOut, nil
}

func CreateJunitReport(testlog string, reportName string) error {
	// Read input
	report, _ := parser.Parse(strings.NewReader(testlog), "")

	// Write xml
	currDir, _ := os.Getwd()

	outfile, err := os.Create(reportName)
	if err != nil {
		console.Errorf("Failed to create report.xml, err %s\n", err)
		return err
	}

	console.Infof("Report file %s created!", currDir+"/"+reportName)

	err = formatter.JUnitReportXML(report, false, "", bufio.NewWriter(outfile))
	if err != nil {
		console.Errorf("Error writing XML: %s\n", err)

	}
	return nil
}
