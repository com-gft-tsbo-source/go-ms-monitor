package msmonitor

import (
	"flag"
	"time"

	"com.gft.tsbo-training.src.go/common/device/devicedb"
	"com.gft.tsbo-training.src.go/common/device/util/devicemap"
	"com.gft.tsbo-training.src.go/common/device/util/devicenode"
	"com.gft.tsbo-training.src.go/common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// MsMonitor
// ###########################################################################
// ###########################################################################

// MsMonitor Encapsulates the ms-monitor data
type MsMonitor struct {
	microservice.MicroService
	*DevicesConfiguration
	devicemap.DeviceMap
	starttime    time.Time
	DBConnection devicedb.IConnection
}

// ---------------------------------------------------------------------------

type imsMonitor interface {
	microservice.IMicroService
	devicemap.IDeviceMap
	GetStarttime() time.Time
	LoadDevice(connectString string) (*devicenode.DeviceNode, error)
	UpdateDevices() error
	StoreDevices() error
}

// ###########################################################################

// InitMsMonitorFromArgs
func InitFromArgs(ms *MsMonitor, args []string, flagset *flag.FlagSet) *MsMonitor {
	var db devicedb.IConnection
	var cfg Configuration
	InitConfigurationFromArgs(&cfg, args, flagset)
	ms.DevicesConfiguration = &cfg.DevicesConfiguration
	microservice.Init(&ms.MicroService, &cfg.Configuration, nil)
	devicemap.InitDeviceMap(&ms.DeviceMap)

	if len(ms.GetDBName()) > 0 {
		db = devicedb.NewDatabase(ms.GetDBName(), "measurements", ms.GetName())
		_, isNil := db.(*devicedb.NilConnection)

		if !isNil {
			ms.GetLogger().Printf("Got database configuration '%s'.\n", ms.GetDBName())
		} else {
			ms.GetLogger().Printf("Bad database configuration '%s', ignoring it.\n", ms.GetDBName())
			db = nil
		}
	} else {
		ms.GetLogger().Println("No database configured.")
	}
	ms.DBConnection = db

	if len(cfg.Devices) >= 1 {
		for i := 0; i < len(cfg.Devices); i++ {
			node, error := ms.LoadDevice(cfg.Devices[i])
			if error != nil {
				ms.GetLogger().Println(error.Error())
				continue
			}
			ms.GetLogger().Printf("Created device '%s'.", node.GetDeviceAddress())
		}
	}

	handlerMonitor := ms.DefaultHandler()
	handlerMonitor.Get = ms.httpGetMonitor
	handlerMonitor.Post = ms.httpPostMonitor
	handlerMonitor.Delete = ms.httpDeleteMonitor
	ms.AddHandler("/monitor", handlerMonitor)
	return ms
}

// ###########################################################################

// GetStarttime gives the starttime of the process
func (ms *MsMonitor) GetStarttime() time.Time { return ms.starttime }

// ---------------------------------------------------------------------------

// LoadDevice ...
func (ms *MsMonitor) LoadDevice(url string) (*devicenode.DeviceNode, error) {

	deviceInfo, err := devicenode.FromURL(url)

	if err != nil {
		return nil, err
	}

	node, err := ms.Add(deviceInfo)

	if err != nil {
		return nil, err
	}

	if ms.DBConnection != nil {

		err = ms.DBConnection.AddDevice(node)

		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

// ---------------------------------------------------------------------------

// UpdateDevices ...
func (ms *MsMonitor) UpdateDevices() error {

	var unregister []string

	for key, node := range ms.ByDeviceAddress {

		error := node.Update()

		if error != nil {
			ms.GetLogger().Printf("Failed to update '%s'. Unregistering.", key)
			unregister = append(unregister, key)
		} else {
			ms.GetLogger().Printf("Got value for '%s': '%s' at '%s'.", key, node.Formatted, node.Stamp.Format("2006/01/02 15:04:05.000"))
		}
	}

	for idx := range unregister {
		key := unregister[idx]
		ms.GetLogger().Printf("Deleting unreachable '%s'.", key)
		ms.Delete(key)
	}

	return nil
}

// ---------------------------------------------------------------------------

// StoreDevices ...
func (ms *MsMonitor) StoreDevices() error {

	var unregister []string

	if ms.DBConnection == nil {
		return nil
	}

	if !ms.DBConnection.IsOpen() {
		return nil
	}

	for key, node := range ms.ByDeviceAddress {

		var err error

		if node.IsNew() {
			err = ms.DBConnection.AddDevice(node)

			if err != nil {
				unregister = append(unregister, key)
				continue
			}
			node.ClearNew()
		}

		err = ms.DBConnection.Update(node)

		if err != nil {
			ms.GetLogger().Printf("Failed to store '%s'! Database connection may be broken.", key)
			ms.DBConnection.Close()
			return err
		}
	}

	for idx := range unregister {
		key := unregister[idx]
		ms.GetLogger().Printf("Deleting unreachable '%s'.", key)
		ms.Delete(key)
	}

	return nil
}
