package antwar

import (
	"exp/gui"
	"exp/gui/x11"
	"container/vector"
	"time"
	"sync"
)

type GUI struct {
	window gui.Window
	queue vector.Vector
	board *Board
	mutex sync.Mutex
}

func (gui *GUI) Update(pos Pos) {
	gui.mutex.Lock()
	gui.queue.Push(pos)
	gui.mutex.Unlock()
}

func (s *GUI) Close() {
	<- s.window.EventChan()
	s.window.Close()
}

func (s *GUI) Flush() {
	s.window.FlushImage();
}

func (gui *GUI) StartLoop() {
	go func() {
		for {
			gui.mutex.Lock()
			updateCount := gui.queue.Len()
			for i := 0; i < updateCount; i++ {
				pos, _ := gui.queue.At(i).(Pos);
				gui.window.Screen().Set(pos.X, pos.Y, gui.board.At(pos).Color())
			}
			gui.queue.Cut(0, updateCount);
			gui.Flush();
			gui.mutex.Unlock()
			time.Sleep(20000000)
		}
	}()
}


func NewGUI(b *Board) *GUI {
	win, err := x11.NewWindow()
	if (err != nil) {
		println(err);
		return nil
	}
	gui := new(GUI)
	gui.window = win
	gui.board = b
	b.OnUpdate = func (p Pos) { gui.Update(p) }
	gui.StartLoop()
	return gui
}