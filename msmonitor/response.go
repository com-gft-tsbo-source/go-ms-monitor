package msmonitor

import (
	"com.gft.tsbo-training.src.go/common/device/util/devicenode"
	"com.gft.tsbo-training.src.go/common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// MsMonitor Response
// ###########################################################################
// ###########################################################################

// MonitorResponse Encapsulates the reploy of ms-monitor
type MonitorResponse struct {
	microservice.Response
	Devices []*devicenode.DeviceNode `json:"devices"`
}

// ###########################################################################

// InitMonitorResponse Constructor of a response of ms-monitor
func InitMonitorResponse(r *MonitorResponse, status string, ms *MsMonitor) {
	microservice.InitResponseFromMicroService(&r.Response, ms, status)
	r.Devices = make([]*devicenode.DeviceNode, len(ms.ByDeviceAddress))
	idx := 0
	for _, value := range ms.ByDeviceAddress {
		r.Devices[idx] = value
		idx = idx + 1
	}
}

// NewMonitorResponse Constructor of a response of ms-monitor
func NewMonitorResponse(status string, ms *MsMonitor) *MonitorResponse {
	var r MonitorResponse
	InitMonitorResponse(&r, status, ms)
	return &r
}
