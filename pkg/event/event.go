package event

type Event interface {
	String() string
}

type Broadcaster interface {
	Broadcast(site string, event Event)
}
