package config

import (
	"os"

	// yaml
	"gopkg.in/yaml.v2"
)

// Fucntion to read a file as text
func ReadFileAsText(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	data := make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Function to read a file as yaml
func ReadFileAsYaml(path string, v interface{}) error {
	data, err := ReadFileAsText(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal([]byte(data), v)
}

/*
	This is the content of config file

healerdb:
    connstr: "mongodb://localhost:27017"
    conncreds:
        username:
        password:
    dbs:
        - name: "enum"
          target_based: true
        - name: "vuln"
          target_based: true
        - name: "watch"
          target_based: true
        - name: "notifio"
          target_based: false
        - name: "report"
          target_based: true
        - name: "schedule"
          target_based: true
        - name: "ca"
          target_based: true
        - name: "web"
          target_based: false
        - name: "creds"
          target_based: false
        - name: "modules_api"
          target_based: false
        - name: "worker"
          target_based: false
        - name: "log"
          target_based: true

Now we should define a Config type based on the above config file
*/

type Config struct {
	HealerDB struct {
		Connstr   string `yaml:"connstr"`
		Conncreds struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"conncreds"`
		Dbs []struct {
			Name        string `yaml:"name"`
			TargetBased bool   `yaml:"target_based"`
		} `yaml:"dbs"`
	} `yaml:"healerdb"`
}

// Function to read the config file and return a Config type
func ReadConfig() (*Config, error) {
	config := &Config{}
	err := ReadFileAsYaml("/ptv/healer/healerdb/config/config.yaml", config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Function GetConnStr to read the connection string from the config file
func GetConnStr() (string, error) {
	config, err := ReadConfig()
	if err != nil {
		return "", err
	}
	return config.HealerDB.Connstr, nil
}

// Function GetConnCreds to read the connection credentials from the config file
func GetConnCreds() (string, string, error) {
	config, err := ReadConfig()
	if err != nil {
		return "", "", err
	}
	return config.HealerDB.Conncreds.Username, config.HealerDB.Conncreds.Password, nil
}

// Function GetAllDatabases : to Read all the databases from the config file and return a slice of struct of them
func GetDatabases() ([]struct {
	Name        string `yaml:"name"`
	TargetBased bool   `yaml:"target_based"`
}, error) {
	config, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	var dbs []struct {
		Name        string `yaml:"name"`
		TargetBased bool   `yaml:"target_based"`
	}
	dbs = append(dbs, config.HealerDB.Dbs...)
	return dbs, nil
}

// Function GetDbs to read the database names from the config file
func GetDatabasesNames() ([]string, error) {
	dbs_names := []string{}
	dbs, err := GetDatabases()
	if err != nil {
		return nil, err
	}
	for _, db := range dbs {
		dbs_names = append(dbs_names, db.Name)
	}
	return dbs_names, nil
}
