package symphony

import "os"

type TaskConfigs struct {
	Filenames []string
}

func NewTaskConfigs(directoryName string) (*TaskConfigs, error) {

	p := new(TaskConfigs)

	f, err := os.Open(directoryName)
	if err != nil {
		return p, err
	}

	fileInfo, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return p, err
	}

	for _, file := range fileInfo {
		p.Filenames = append(p.Filenames, file.Name())
	}

	return p, nil

}

func (tcs *TaskConfigs) EnumerateFilenames() []string {
	return tcs.Filenames
}
