//
// Rover - Landing zone and launchpad actions
// * Encapsulation of actions
// * Ben C, May 2021
//

package landingzone

type Action int

var actionEnum = []string{"init", "plan", "apply", "run", "destroy", "test", "format", "validate"}
var descriptionEnum = []string{
	"Perform a terraform init and no other action",
	"Perform a terraform plan",
	"Perform a terraform plan & apply",
	"Perform a terraform plan, apply & test",
	"Perform a terraform destroy",
	"Run a test using terratest",
	"Perform a terraform fmt check",
	"Perform a terraform validate",
}

const (
	// ActionInit carries out a init operation and exits
	ActionInit Action = iota
	// ActionPlan carries out a plan operation
	ActionPlan Action = iota
	// ActionApply carries out a plan AND apply
	ActionApply Action = iota
	// ActionRun carries out a plan, apply + test
	ActionRun Action = iota
	// ActionDestroy carries out a destroy operation
	ActionDestroy Action = iota
	// ActionTest carries out a test operation
	ActionTest Action = iota
	// ActionFormat carries out a fmt operation
	ActionFormat Action = iota
	// ActionValidate carries out a vaildate operation
	ActionValidate Action = iota
)

// NewAction returns an Action type from a string
// func NewAction(actionString string) (Action, error) {
// 	switch strings.ToLower(actionString) {
// 	case ActionInit.String():
// 		return ActionInit, nil
// 	case ActionPlan.String():
// 		return ActionPlan, nil
// 	case ActionRun.String():
// 		return ActionRun, nil
// 	case ActionDestroy.String():
// 		return ActionDestroy, nil
// 	default:
// 		return 0, errors.New("action is not valid")
// 	}
// }

func (a Action) Name() string {
	return actionEnum[a]
}

func (a Action) Description() string {
	return descriptionEnum[a]
}

// func AddActionFlag(cmd *cobra.Command) {
// 	cmd.Flags().StringP("action", "a", "init", fmt.Sprintf("Action to run, one of %v", actionEnum))
// }
