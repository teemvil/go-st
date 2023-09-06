package structures

import "time"

type ManagementMessage struct {
	DeviceName       string    `json: "name"`
	Itemid           string    `json: "itemid"`
	Message          string    `json: "message"`
	Event            string    `json: "event"`
	Time             time.Time `json: "time"`
	Jwt              string    `json: "jwt"`
	SensorName       string    `json: "sensorName"`
	SensorHostDevice string    `json: "sensorHostDevice"`
	SensorChannel    string    `json: "sensorChannel"`
	Misc             string    `json: "misc"`
}

type SensorFile struct {
	Name        string `json: "name"`
	HostDevice  string `json: "hostdevice"`
	MQTTchannel string `json: "mqttchannel"`
}

type DeviceFile struct {
	Name   string `json: "name"`
	Itemid string `json: "itemid"`
}
