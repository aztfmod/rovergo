//
// Rover - Unit tests for utility methods
//

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Read_Yaml_File_Without_Extension(t *testing.T) {

	// Arrange
	pwd := os.Getenv("PWD")
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/examples/custom_commands/commands"

	// Act
	b, err := ReadYamlFile(fileName)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

}

func Test_Read_Yaml_File_With_Yml_Extension(t *testing.T) {

	// Arrange
	pwd := os.Getenv("PWD")
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/examples/custom_commands/commands.yml"

	// Act
	b, err := ReadYamlFile(fileName)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

}

func Test_Read_Yaml_File_With_Yaml_Extension(t *testing.T) {

	// Arrange
	pwd := os.Getenv("PWD")
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/examples/custom_commands/commands.yaml"

	// Act
	b, err := ReadYamlFile(fileName)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

}

func Test_Read_Not_Exist_File(t *testing.T) {

	// Arrange
	pwd := os.Getenv("PWD")
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/examples/custom_commands/commands.eng"

	// Act
	_, err := ReadYamlFile(fileName)

	// Assert
	assert.Error(t, err)

}
