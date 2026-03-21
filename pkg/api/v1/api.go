package api

// Client defines the client interface for the API.
type Client interface {
	Status() (*StatusReply, error)
}

// Reply contains standard fields for generic API replies.
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Returned on status requests.
type StatusReply struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime,omitempty"`
	Version string `json:"version,omitempty"`
}
