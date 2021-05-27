//
// Rover - main entry point
// * Simply enters into the cobra root command and passes version
// * Ben C, May 2021
//

package main

import (
	"github.com/aztfmod/rover/cmd"
)

// Have to put version here or ldflags can't set it ¯\_(ツ)_/¯
var version = "0.0.0"

func main() {
	cmd.Execute(version)
}
