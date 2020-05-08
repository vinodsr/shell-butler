package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	types "github.com/vinodsr/shell-butler/lib/types"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//WriteConfig Write the config to a file
func WriteConfig(config types.ConfigData) {
	dir, _ := os.UserHomeDir()
	configFile := dir + "/.butler/settings.json"
	initStructJSON, _ := json.MarshalIndent(config, "", "    ")
	_, err := os.Stat(dir + "/.butler")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dir+"/.butler", 0755)
		check(errDir)

	}
	err = ioutil.WriteFile(configFile, initStructJSON, 0644)
	check(err)
}
