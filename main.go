package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	addcli "github.com/vinodsr/shell-butler/lib/cli/add"
	"github.com/vinodsr/shell-butler/lib/cli/help"
	config "github.com/vinodsr/shell-butler/lib/config"
	runtime "github.com/vinodsr/shell-butler/lib/runtime"
	types "github.com/vinodsr/shell-butler/lib/types"
	layout "github.com/vinodsr/shell-butler/lib/ui"
)

//Main program
func main() {
	var configData types.ConfigData

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	dir, _ = os.UserHomeDir()
	configFile := dir + "/.butler/settings.json"
	dat, err := ioutil.ReadFile(configFile)
	if err != nil {
		initStruct := types.ConfigData{
			Commands: []types.Command{{
				Context: "Help",
				Program: "echo Edit " + configFile + " to add more commands",
			}},
		}
		//fmt.Printf("%+v\n", initStruct)
		//os.Exit(0)
		config.WriteConfig(initStruct)
		configData = initStruct
	}

	json.Unmarshal([]byte(string(dat)), &configData)
	//fmt.Printf("Data %+v\n", commands)
	//fmt.Print(commands.Commands[0].Context)

	// initialise the list datasource .

	// Initialize the runtime
	var rt *runtime.Runtime = runtime.GetRunTime()
	rt.Init(configData)

	if len(os.Args) > 1 {
		argCommand := os.Args[1]
		if argCommand == "add" {
			addcli.Execute()
			os.Exit(1)
		} else if argCommand == "help" || argCommand == "--help" {
			help.Execute()
			os.Exit(1)
		}

	}

	layout.Render()
	//os.Exit(1)

}
