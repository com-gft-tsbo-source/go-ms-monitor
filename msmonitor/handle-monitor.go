package msmonitor

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"com.gft.tsbo-training.src.go/common/device/util/devicemap"
)

// ---------------------------------------------------------------------------

func (ms *MsMonitor) httpGetMonitor(w http.ResponseWriter, r *http.Request) (int, contentLen int, msg string) {
	status := http.StatusOK
	response := NewMonitorResponse("OK", ms)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("cid", ms.GetName())
	w.Header().Set("version", ms.GetVersion())
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, "Reported monitored devices."
}

// ---------------------------------------------------------------------------

func (ms *MsMonitor) httpPostMonitor(w http.ResponseWriter, r *http.Request) (int, contentLen int, msg string) {

	status := http.StatusCreated
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		msg = fmt.Sprintf("Failed to read request body, error was '%s'!", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return http.StatusInternalServerError, 0, msg
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "Empty device hostname.", http.StatusBadRequest)
		return http.StatusInternalServerError, 0, "Empty device hostname."
	}

	connectstring := string(body[:len(body)])
	node, err := ms.LoadDevice(connectstring)

	if err != nil {
		msg = fmt.Sprintf("Failed to load device from '%s', error was '%s'!", connectstring, err.Error())
		switch err.(type) {
		case *devicemap.OpErrorDuplicate, devicemap.OpErrorDuplicate:
			http.Error(w, msg, http.StatusConflict)
			return http.StatusConflict, 0, msg
		}
		http.Error(w, msg, http.StatusBadRequest)
		return http.StatusInternalServerError, 0, msg
	}

	msg = fmt.Sprintf("Device '%s' created.", node.GetDeviceAddress())
	response := NewMonitorResponse(msg, ms)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("cid", ms.GetName())
	w.Header().Set("version", ms.GetVersion())
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, msg
}

// ---------------------------------------------------------------------------

func (ms *MsMonitor) httpDeleteMonitor(w http.ResponseWriter, r *http.Request) (int, contentLen int, msg string) {
	status := http.StatusOK
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		msg = fmt.Sprintf("Failed to read request body, error was '%s'!", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return http.StatusInternalServerError, 0, "Failed to read body."
	}

	defer r.Body.Close()
	if len(body) == 0 {
		http.Error(w, "Empty device hostname.", http.StatusBadRequest)
		return http.StatusInternalServerError, 0, "Empty device hostname."
	}

	connectstring := string(body[:len(body)])
	deviceAddress, _, err := net.SplitHostPort(connectstring)

	if err != nil {
		deviceAddress = connectstring
	}
	node, err := ms.Delete(deviceAddress)

	if err != nil {
		msg = fmt.Sprintf("Cannot delete device '%s', error was '%s'!", deviceAddress, err.Error())
		switch err.(type) {
		case *devicemap.OpErrorNotFound, devicemap.OpErrorNotFound:
			http.Error(w, msg, http.StatusNotFound)
			return http.StatusNotFound, 0, msg
		}
		http.Error(w, msg, http.StatusBadRequest)
		return http.StatusInternalServerError, 0, msg
	}

	msg = fmt.Sprintf("Device '%s' deleted.", node.GetDeviceAddress())
	response := NewMonitorResponse(msg, ms)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("cid", ms.GetName())
	w.Header().Set("version", ms.GetVersion())
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, msg
}
