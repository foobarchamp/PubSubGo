package node

import (
	"math/rand/v2"
)

type MsgType int

const (
	Closing MsgType = iota
	subscribe
	unsubscribe
	Number
	TypeResponse
)

const (
	BUFF_LEN  = 1_000_000
	MSG_COUNT = 10_000_000
)

var (
//	Msgpool = sync.Pool{
//		New: func() interface{} {
//			return new(Message)
//		},
//	}
)

type Message struct {
	Type    MsgType
	Payload interface{}
}

type Node struct {
	id          int
	Name        string
	Subscribers map[int]*Node
	Subscribed  map[int]*Node
	inQ         chan *Message
}

func (main *Node) addSubscriber(subscriber *Node) {
	main.GetNode().Subscribers[subscriber.id] = subscriber
	subscriber.Subscribed[main.id] = main.GetNode()
}

func (subscriber *Node) RequestSubscriptionTo(main PubSubNode) {
	main.GetNode().inQ <- &Message{
		Type:    subscribe,
		Payload: subscriber,
	}
}

func (subscriber *Node) RequestUnSubscriptionFrom(main *Node) {
	main.inQ <- &Message{
		Type:    unsubscribe,
		Payload: subscriber,
	}
}

func (main *Node) removeSubscriber(subscriber *Node) {
	delete(main.Subscribers, subscriber.id)
	delete(subscriber.Subscribed, main.id)
}

func (n *Node) GetNode() *Node {
	return n
}

func (n *Node) Signal(msg *Message) {
	for _, subscriber := range n.Subscribers {
		subscriber.inQ <- msg
	}
}

func NewNode(name string) *Node {
	return &Node{
		id:          rand.Int(),
		Name:        name,
		Subscribers: make(map[int]*Node),
		Subscribed:  make(map[int]*Node),
		inQ:         make(chan *Message, BUFF_LEN),
	}
}

type PubSubNode interface {
	GetNode() *Node
	NextMessage(msg *Message)
	RequestSubscriptionTo(main PubSubNode)
}

func HandleInQ(psnode PubSubNode) {
	for msg := range psnode.GetNode().inQ {
		switch msg.Type {
		case subscribe:
			psnode.GetNode().addSubscriber(msg.Payload.(*Node))
		case unsubscribe:
			psnode.GetNode().removeSubscriber(msg.Payload.(*Node))
		default:
			psnode.NextMessage(msg)
		}
	}
}

func Run(psnode PubSubNode) {
	go HandleInQ(psnode)
}
