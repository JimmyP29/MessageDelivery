package messaging

type hub struct {
	clients  map[uint64]*client
	messages chan message
}
