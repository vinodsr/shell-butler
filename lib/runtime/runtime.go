package runtime

import (
	"sync"

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
func (rt *Runtime) Init(configData types.ConfigData, commandMap map[string]string, contextDataSource []string) {
	rt.configData = configData
	rt.commandMap = commandMap
	rt.initialized = true
	rt.contextDataSource = contextDataSource
}
