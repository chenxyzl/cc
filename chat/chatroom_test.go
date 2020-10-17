package chat

import (
	"fmt"
	"testing"
	"time"
)

func TestJoin(t *testing.T) {

	room := NewRoom()

	a := room.Join("a", 1)
	go func() {
		for v := range a.Pipe {
			fmt.Println(v)
		}
	}()

	b := room.Join("b", 2)
	go func() {
		for v := range b.Pipe {
			fmt.Println(v)
		}
	}()
	a.Say("aa")
	b.Say("bb")

	a.Leave()
	b.Leave()
	time.Sleep(1 * time.Second)

}
