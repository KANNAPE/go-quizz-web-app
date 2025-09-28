package lobby

type Lobby struct {
	ID       string
	HostName string
	// Participants []string  // extend later
}

type Store interface {
	Create(host string) (Lobby, error)
	Get(id string) (Lobby, bool)
}
