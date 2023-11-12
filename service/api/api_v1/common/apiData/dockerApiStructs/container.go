package dockerApiStructs

type PostType string

const (
	Tcp  PostType = "tcp"
	Udp  PostType = "udp"
	Sctp PostType = "sctp"
)

type ContainerPort struct {
	IP          string   `json:"ip"`
	PrivatePort int      `json:"privatePort"`
	PublicPort  int      `json:"publicPort"`
	Type        PostType `json:"type"`
}

type ContainerInfoResp struct {
	Name       string          `json:"name"`
	Image      string          `json:"image"`
	ImageID    string          `json:"imageID,omitempty"`
	Command    string          `json:"command,omitempty"`
	CreateTime string          `json:"createTime"`
	Ports      []ContainerPort `json:"ports"`
	State      string          `json:"state"`
}
