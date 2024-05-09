package dockerApiStructs

import "time"

type PostType string

const (
	Tcp  PostType = "tcp"
	Udp  PostType = "udp"
	Sctp PostType = "sctp"
)

type ContainerPort struct {
	IP          string   `json:"ip"`
	PrivatePort uint16   `json:"privatePort"`
	PublicPort  uint16   `json:"publicPort"`
	Type        PostType `json:"type"`
}

type ContainerInfo struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Image      string          `json:"image"`
	ImageID    string          `json:"imageID,omitempty"`
	Command    string          `json:"command,omitempty"`
	CreateTime string          `json:"createTime"`
	Ports      []ContainerPort `json:"ports"`
	State      string          `json:"state"`
	Status     string          `json:"status"`
}
type ContainerInfoResp struct {
	ContainerInfo
}

type ExecSwitchActionByContainerIDReq struct {
	ID     string `json:"id"`
	Action bool   `json:"action"`
}

type ContainerStats struct {
	CPUPercent   float64 `json:"cpuPercent"`
	Memory       float64 `json:"memory"`
	Cache        float64 `json:"cache"`
	IORead       float64 `json:"ioRead"`
	IOWrite      float64 `json:"ioWrite"`
	NetworkRX    float64 `json:"networkRX"`
	NetworkTX    float64 `json:"networkTX"`
	RunStatus    string  `json:"runStatus"`
	ErrorMassage string  `json:"errorMessage"`

	ShotTime time.Time `json:"shotTime"`
}
