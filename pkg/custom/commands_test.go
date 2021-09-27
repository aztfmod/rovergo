//
// Rover - Unit tests for custom commands and group commands
//

package custom

import (
	"testing"

	"github.com/aztfmod/rover/pkg/console"
	"github.com/stretchr/testify/assert"
)

func Test_Commands_Yaml_Doesnt_Exists(t *testing.T) {

	console.DebugEnabled = true

	actions, err := LoadCustomCommandsAndGroups()

	assert.Error(t, err)

	assert.Empty(t, actions)

}
