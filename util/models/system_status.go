package models

type SystemStatus struct {
	Timestamp     int64  `json:"timestamp"`
	ValidVersion  bool   `json:"valid_version"`
	SystemRunnnig bool   `json:"system_running"`
	Message       string `json:"message"`
}
