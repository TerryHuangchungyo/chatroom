package wsservice

type Client struct {
	Id   uint32
	Name string
	Hubs map[uint32]bool
}
