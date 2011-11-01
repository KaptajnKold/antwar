package antwar

import (
	"exp/gui"
	"exp/gui/x11"
	"container/vector"
	"time"
)

type GUI struct {
	window gui.Window
	queue vector.Vector
	board *Board
}

func (gui *GUI) Update(pos Pos) {
	gui.queue.Push(pos)
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
			updateCount := gui.queue.Len()
			for i := 0; i < updateCount; i++ {
				pos, _ := gui.queue.At(i).(Pos);
				gui.window.Screen().Set(pos.X, pos.Y, gui.board.At(pos).Color())
			}
			gui.queue.Cut(0, updateCount);
			gui.Flush();
			time.Sleep(200000000)
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