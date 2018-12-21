package internal

import "testing"

func TestLoadConfigurationFile(t *testing.T) {
	var configuration Configuration
	jsonData, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonData, &configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
}