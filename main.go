package main

import (
	"PubSubGo/node"
	"fmt"
	"sync"
	"time"
)

var (
	now time.Time
)

type Doubler struct {
	*node.Node
	count int
	wg    *sync.WaitGroup
}

func NewDoubler(name string) *Doubler {
	d := &Doubler{node.NewNode(name), 0, &sync.WaitGroup{}}
	d.wg.Add(1)
	return d
}

func (d *Doubler) NextMessage(msg *node.Message) {
	d.count++
	if d.count == node.MSG_COUNT {
		diff := time.Since(now).Milliseconds()
		fmt.Printf("%s processed %d messages in %d ms\n", d.Name, node.MSG_COUNT, diff)
		d.wg.Done()
		//node.Stop(d)
	}
	//d.Signal(&node.Message{Type: msg.Type, Payload: 2 * msg.Payload.(int)})
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
		msg := node.Msgpool.Get().(*node.Message)
		msg.Payload = i
		nsp.Signal(msg)
	}
}

func main() {
	nsp := NewNSP("NSP")
	node.Run(nsp)
	d := NewDoubler("Doubler")
	node.Run(d)
	d.SubscribeTo(nsp)
	time.Sleep(1 * time.Second)
	now = time.Now()
	nsp.start()
	d.wg.Wait()
	//d := 0
	//for range node.MSG_COUNT {
	//	d++
	//}
	//fmt.Println(time.Since(now).Microseconds())
}
