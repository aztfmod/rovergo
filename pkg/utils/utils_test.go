//go:build unit
// +build unit

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Read_Yaml_File_Without_Extension(t *testing.T) {

	// Arrange
	pwd, _ := os.Getwd()
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/examples/custom_commands/commands"

	// Act
	fileContent, _, err := ReadYamlFile(fileName)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, fileContent)

}

func Test_Read_Yaml_File_With_Yml_Extension(t *testing.T) {

	// Arrange
	pwd, _ := os.Getwd()
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/examples/custom_commands/commands.yml"

	// Act
	fileContent, _, err := ReadYamlFile(fileName)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, fileContent)

}

func Test_Read_Yaml_File_With_Yaml_Extension(t *testing.T) {

	// Arrange
	pwd, _ := os.Getwd()
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/test/testdata/custom_commands/_default.yaml"

	// Act
	fileContent, _, err := ReadYamlFile(fileName)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, fileContent)

}

func Test_Read_Not_Exist_File(t *testing.T) {

	// Arrange
	pwd, _ := os.Getwd()
	oneUp := filepath.Dir(pwd)
	projectRoot := filepath.Dir(oneUp)
	fileName := projectRoot + "/examples/custom_commands/commands.eng"

	// Act
	fileContent, fileName, err := ReadYamlFile(fileName)

	// Assert
	assert.EqualError(t, err, "file extension must be .yaml or .yml")
	assert.Nil(t, fileContent)
	assert.Equal(t, "commands.eng", fileName)
}
