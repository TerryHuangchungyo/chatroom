package wsservice

type Hub struct {
	Id         uint32
	Name       string
	Clients	   map[uint32]bool
}

