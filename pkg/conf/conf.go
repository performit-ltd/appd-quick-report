package conf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Define the YAML conf struct
type Conf struct {
	Stats StatsConf `yaml:"stats"`
}
type StatsConf []struct {
	Name    string     `yaml:"name"`
	Url     string     `yaml:"url"`
	Client  string     `yaml:"client"`
	Secret  string     `yaml:"secret"`
	Account string     `yaml:"account"`
	Auth    string     `yaml:"auth"`
	Report  ReportConf `yaml:"report"`
}
type ReportConf struct {
	Name        string     `yaml:"name"`
	Subtitle    string     `yaml:"subtitle"`
	Timerange   string     `yaml:"timerange"`
	Scope       string     `yaml:"scope"`
	Team        string     `yaml:"team"`
	Description string     `yaml:"description"`
	Header      HeaderConf `yaml:"header"`
}
type HeaderConf struct {
	B2 string `yaml:"b2"`
	B3 string `yaml:"b3"`
	B4 string `yaml:"b4"`
	B5 string `yaml:"b5"`
}

func LoadConf() Conf {

	// Set the conf variable to return
	var yamlconf Conf

	// Set output color vars
	var Red = "\033[31m"
	var Reset = "\033[0m"

	// Read the YAML conf file in memory
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Printf("\n%vFailed to read conf file.%v\n\n", Red, Reset)
		log.Fatal(err)
		os.Exit(1)
	}

	// Unmarshal YAML conf into the given struct
	err = yaml.Unmarshal(yamlFile, &yamlconf)
	if err != nil {
		fmt.Printf("\n%vFailed to unmarshal conf file.%v\n\n", Red, Reset)
		log.Fatal(err)
		os.Exit(1)
	}

	// Return the loaded conf struct
	log.Println("YAML config loaded successfully.")

	return yamlconf

}
