package node

import (
	"math/rand/v2"
	"sync"
)

//type MsgType int

//const (
//	//SysMsg MsgType = iota
//	Number
//	TypeResponse
//)

const (
	BUFF_LEN  = 1_000_000
	MSG_COUNT = 10_000_000
)

var (
	Msgpool = sync.Pool{
		New: func() interface{} {
			return new(Message)
		},
	}
)

type Message struct {
	//Type    MsgType
	Payload interface{}
}

type Node struct {
	id          int
	Name        string
	subscribers map[int]*Node
	subscribed  map[int]*Node
	inQ         chan *Message
	outQ        chan *Message
}

func (subscriber *Node) SubscribeTo(main PubSubNode) {
	main.GetNode().subscribers[subscriber.id] = subscriber
	subscriber.subscribed[main.GetNode().id] = main.GetNode()
}

func (subscriber *Node) UnsubscribeFrom(main *Node) {
	delete(main.subscribers, subscriber.id)
	delete(subscriber.subscribed, main.id)
}

func (n *Node) GetNode() *Node {
	return n
}

func (n *Node) Signal(msg *Message) {
	n.outQ <- msg
}

func NewNode(name string) *Node {
	return &Node{
		id:          rand.Int(),
		Name:        name,
		subscribers: make(map[int]*Node),
		subscribed:  make(map[int]*Node),
		inQ:         make(chan *Message, BUFF_LEN),
		outQ:        make(chan *Message, BUFF_LEN),
	}
}

type PubSubNode interface {
	GetNode() *Node
	NextMessage(msg *Message)
	SubscribeTo(main PubSubNode)
	//Run() error
}

func HandleInQ(psnode PubSubNode) {
	for msg := range psnode.GetNode().inQ {
		//fmt.Println("pulling", msg.Payload)
		psnode.NextMessage(msg)
		Msgpool.Put(msg)
	}
}

func HandleOutQ(psnode PubSubNode) {
	for msg := range psnode.GetNode().outQ {
		//if msg.Type == SysMsg {
		//	fmt.Println("//handle ", msg)
		//} else {
		//	// Process message
		//	if len(psnode.GetNode().subscribers) != 0 {
		//	}
		//}
		//fmt.Println("pushing", msg.Payload)
		for _, subscriber := range psnode.GetNode().subscribers {
			subscriber.inQ <- msg
		}
	}
}

func Run(psnode PubSubNode) {
	go HandleOutQ(psnode)
	go HandleInQ(psnode)
}
