//
// Rover - main entry point
// * Checks for dependent executables then
// * Simply enters into the cobra root command and passes version
// * Ben C, May 2021
//

package main

import (
	"github.com/aztfmod/rover/cmd"
	"github.com/aztfmod/rover/pkg/command"
)

func main() {
	command.ValidateDependencies()
	cmd.Execute()
}
