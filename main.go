package main

import (
	"PubSubGo/node"
	"fmt"
	"sync"
	"time"
)

var (
	now time.Time
	wg  *sync.WaitGroup
)

type Doubler struct {
	*node.Node
	count int
}

func NewDoubler(name string) *Doubler {
	d := &Doubler{node.NewNode(name), 0}
	return d
}

func (d *Doubler) NextMessage(msg *node.Message) {
	switch msg.Type {
	case node.Number:
		d.count++
		d.Signal(&node.Message{Type: msg.Type, Payload: 2 * msg.Payload.(int)})
	default:
		panic(fmt.Sprintf("unknown message %v", msg))
	}
}

func (d *Doubler) Suspend() {
	diff := time.Since(now).Milliseconds()
	d.GetNode().Suspend()
	fmt.Println("processed", d.count, "message in", diff, "ms")
	wg.Done()
}

type NSP struct {
	*node.Node
}

func NewNSP(name string) *NSP {
	return &NSP{node.NewNode(name)}
}

func (d *NSP) NextMessage(msg *node.Message) {
	panic("not implemented")
}

func (nsp *NSP) start() {
	fmt.Println("Starting", nsp.Name)
	for i := range node.MSG_COUNT {
		nsp.Signal(&node.Message{
			Type:    node.Number,
			Payload: i,
		})
	}
	nsp.Suspend()
	fmt.Println("Finished", nsp.Name)
}

func main() {
	nsp := NewNSP("NSP")
	node.Run(nsp)
	d1 := NewDoubler("D1")
	node.Run(d1)
	d1.RequestSubscriptionTo(nsp)
	d2 := NewDoubler("D2")
	node.Run(d2)
	d2.RequestSubscriptionTo(nsp)
	d3 := NewDoubler("D3")
	node.Run(d3)
	d3.RequestSubscriptionTo(d1)
	d3.RequestSubscriptionTo(d2)
	wg = &sync.WaitGroup{}
	wg.Add(4)
	time.Sleep(1 * time.Second)
	now = time.Now()
	nsp.start()
	wg.Wait()
}
