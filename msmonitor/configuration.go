package msmonitor

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

type connectStrings []string

func (cs *connectStrings) String() string { return "" }
func (cs *connectStrings) Set(value string) error {
	*cs = append(*cs, value)
	return nil
}

// DevicesConfiguration ...
type DevicesConfiguration struct {
	Devices []string `json:"devices"`
}

// IDevicesConfiguration ...
type IDevicesConfiguration interface {
	GetDevices() *[]string
}

// Configuration ...
type Configuration struct {
	microservice.Configuration
	DevicesConfiguration
}

// IConfiguration ...
type IConfiguration interface {
	microservice.IConfiguration
	IDevicesConfiguration
}

// GetDevices ...
func (cfg DevicesConfiguration) GetDevices() *[]string { return &cfg.Devices }

// ---------------------------------------------------------------------------

// InitConfigurationFromArgs ...
func InitConfigurationFromArgs(cfg *Configuration, args []string, flagset *flag.FlagSet) {

	var csCli connectStrings

	if flagset == nil {
		flagset = flag.NewFlagSet("ms-monitor", flag.PanicOnError)
	}

	flagset.Var(&csCli, "device", "Connectstrings")
	microservice.InitConfigurationFromArgs(&cfg.Configuration, args, flagset)
	if len(csCli) > 0 {
		cfg.Devices = csCli
	}

	if len(cfg.GetConfigurationFile()) > 0 {
		file, err := os.Open(cfg.GetConfigurationFile())

		if err != nil {
			flagset.Usage()
			panic(fmt.Sprintf(fmt.Sprintf("Error: Failed to open onfiguration file '%s'. Error was %s!\n", cfg.GetConfigurationFile(), err.Error())))
		}

		defer file.Close()
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err = decoder.Decode(&configuration)
		if err != nil {
			flagset.Usage()
			panic(fmt.Sprintf(fmt.Sprintf("Error: Failed to parse onfiguration file '%s'. Error was %s!\n", cfg.GetConfigurationFile(), err.Error())))
		}
		if len(cfg.Devices) == 0 {
			cfg.Devices = configuration.Devices
		}
	}

}
