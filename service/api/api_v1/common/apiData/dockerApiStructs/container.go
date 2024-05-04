package dockerApiStructs

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

type ContainerInfoResp struct {
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

type ExecSwitchActionByContainerIDReq struct {
	ID     string `json:"id"`
	Action bool   `json:"action"`
}
