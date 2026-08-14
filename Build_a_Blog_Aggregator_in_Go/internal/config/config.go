package config


import(
	"os"
	"path/filepath"
	"encoding/json"
)


type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}


const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(userHomeDir, configFileName)
	return configFilePath, nil
}


func Read() (Config, error) {
	config := Config{}

	ConfigFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	fileContent, err := os.ReadFile(ConfigFilePath)
	if err !=nil {
		return Config{}, err
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}


func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	ConfigFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(ConfigFilePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}