package main

import (
	"os"
	"time"

	"github.com/com-gft-tsbo-source/go-ms-monitor/msmonitor"
)

// ###########################################################################
// ###########################################################################
// MAIN
// ###########################################################################
// ###########################################################################

// ---------------------------------------------------------------------------

func main() {
	var ms msmonitor.MsMonitor
	msmonitor.InitFromArgs(&ms, os.Args, nil)
	go func() {

		if ms.DBConnection == nil {
			return
		}
		for ever := true; ever; ever = true {
			var err error

			if ms.DBConnection.IsOpen() {
				goto sleep
			}
			ms.GetLogger().Println("Trying to connect to database.")
			err = ms.DBConnection.Open()

			if err != nil {
				ms.GetLogger().Println(err.Error())
				goto sleep
			}

			ms.GetLogger().Println("Successfully connected to database.")

		sleep:
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		if ms.DBConnection != nil {
			ms.DBConnection.Open()
		}
		for ever := true; ever; ever = true {
			time.Sleep(1 * time.Second)
			ms.UpdateDevices()
			ms.StoreDevices()
		}
		if ms.DBConnection != nil {
			ms.DBConnection.Close()
		}
	}()

	ms.Run()

}
