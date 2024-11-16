package angus

var (
	clients = make(map[string]*Client)
)

type Client struct {
	events map[string]func()
}
