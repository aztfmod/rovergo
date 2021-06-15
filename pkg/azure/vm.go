//
// Rover - Azure Virtual Machine
// * Greg O, June 2021
//

package azure

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Compute struct {
	Name              string
	ResourceGroupName string
}

type Metadata struct {
	Compute Compute
}

func VMInstanceMetadataService() *Metadata {

	client := http.Client{}

	req, err := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
	if err != nil {
		return nil
	}

	req.Header.Add("metadata", "true")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
	}

	meta := &Metadata{}
	err = json.Unmarshal([]byte(bodyBytes), meta)
	if err != nil {
		return nil
	}

	return meta
}
