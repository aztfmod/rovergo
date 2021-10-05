package landingzone

type Action interface {
	Execute(o *Options) error
	GetType() string
	GetName() string
	GetDescription() string
}

const (
	BuiltinCommand = "Builtin"
	CustomCommand  = "Custom"
	GroupCommand   = "Group"
)

type ActionBase struct {
	Name        string
	Type        string
	Description string
}

type TerraformAction struct {
	ActionBase
	launchPadStorageID string
}

func (ab ActionBase) GetName() string {
	return ab.Name
}

func (ab ActionBase) GetType() string {
	return ab.Type
}

func (ab ActionBase) GetDescription() string {
	return ab.Description
}
