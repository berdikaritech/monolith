package micro

type Client interface {
	Name() string
}

type ClientSet struct {
	clients map[string]Client
}

func NewClientSet(clients ...Client) *ClientSet {
	c := &ClientSet{make(map[string]Client)}
	c.AddClient(clients...)

	return c
}

func (c *ClientSet) AddClient(clients ...Client) {
	for _, client := range clients {
		c.clients[client.Name()] = client
	}
}

func (c *ClientSet) Get(key string) Client {
	return c.clients[key]
}

func (c *ClientSet) Keys() []string {
	keys := []string{}

	for key := range c.clients {
		keys = append(keys, key)
	}

	return keys
}

func (c *ClientSet) Length() int {
	return len(c.clients)
}
