package event

type Event interface {
	String() string
}

type EventBroadcaster interface {
	Broadcast(site string, event Event)
}
