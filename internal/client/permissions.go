package client

type Permission struct {
	Admin bool `json:"admin"`
	Pull  bool `json:"pull"`
	Push  bool `json:"push"`
}
