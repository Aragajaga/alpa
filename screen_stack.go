package main

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScreenStackElem struct {
	screen IScreen
	next   *ScreenStackElem
}

type ScreenStack struct {
	head  *ScreenStackElem
	mutex sync.Mutex
	Size  int
}

func (ss *ScreenStack) Push(s IScreen) {
	ss.mutex.Lock()
	el := new(ScreenStackElem)
	el.next = ss.head
	el.screen = s

	ss.head = el
	ss.Size++
	ss.mutex.Unlock()
}

// Remove and return screen interface on the peak of the stack
//
// If there is no screen at the peak of stack, the function will return a nil value
func (ss *ScreenStack) Pop() IScreen {
	if ss.head == nil {
		return nil
	}

	ss.mutex.Lock()
	el := ss.head
	ss.head = ss.head.next
	ss.Size--
	ss.mutex.Unlock()
	return el.screen
}

// Reverse-iterate screen and call Draw() method for each
//
// This function is altering stack in between by double-reversing, so be careful
func (ss *ScreenStack) Draw(screen *ebiten.Image) {

	if ss.head == nil {
		return
	}

	ss.mutex.Lock()
	{
		var prev *ScreenStackElem = nil
		var current *ScreenStackElem = ss.head
		var next *ScreenStackElem = nil

		for current != nil {
			next = current.next
			current.next = prev
			prev = current
			current = next
		}
		ss.head = prev
	}

	for i := ss.head; i != nil; i = i.next {
		i.screen.Draw(screen)
	}

	{
		var prev *ScreenStackElem = nil
		var current *ScreenStackElem = ss.head
		var next *ScreenStackElem = nil

		for current != nil {
			next = current.next
			current.next = prev
			prev = current
			current = next
		}
		ss.head = prev
	}

	ss.mutex.Unlock()
}

func NewScreenStack() *ScreenStack {
	ss := new(ScreenStack)
	return ss
}
