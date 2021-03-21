package runtime

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/vinodsr/shell-butler/lib/config"
	types "github.com/vinodsr/shell-butler/lib/types"
)

//Runtime type
type Runtime struct {
	commandMap        map[string]string
	configData        types.ConfigData
	initialized       bool
	contextDataSource []string
}

var once sync.Once

var singleton *Runtime

// GetRunTime Gets the runtime
func GetRunTime() *Runtime {
	once.Do(func() {
		singleton = &Runtime{}
	})
	return singleton
}

func (rt *Runtime) isInitialized() bool {
	return rt.initialized
}

// GetCommandMap :
func (rt *Runtime) GetCommandMap() map[string]string {
	return rt.commandMap
}

//GetConfig :
func (rt *Runtime) GetConfig() types.ConfigData {
	return rt.configData
}

// GetContextDS :
func (rt *Runtime) GetContextDS() []string {
	return rt.contextDataSource
}

// Init initializes
func (rt *Runtime) Init() {

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
		rt.configData = initStruct
	}

	json.Unmarshal([]byte(string(dat)), &rt.configData)

	rt.initialized = true
	rt.commandMap = make(map[string]string)
	rt.contextDataSource = nil
	for _, command := range rt.configData.Commands {
		// splits := strings.Split(command.Context, ":")
		rt.contextDataSource = append(rt.contextDataSource, command.Context)
		rt.commandMap[command.Context] = command.Program
	}
}
