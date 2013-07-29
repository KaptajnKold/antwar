package antwar

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"image"
	"sync"
	"time"
)

type GUI struct {
	board  *Board
	mutex  sync.Mutex
	canvas *xgraphics.Image
	queue  map[Pos]bool
	win    *xwindow.Window
}

func (gui *GUI) Update(pos Pos) {
	gui.mutex.Lock()
	gui.queue[pos] = true
	gui.mutex.Unlock()
}

func (gui *GUI) StartLoop() {
	go func() {
		for {
			gui.mutex.Lock()
			for pos, _ := range gui.queue {
				gui.canvas.Set(pos.X, pos.Y, gui.board.At(pos).Color())
				delete(gui.queue, pos)
			}
			gui.mutex.Unlock()
			gui.canvas.XDraw()
			gui.canvas.XPaint(gui.win.Id)
			time.Sleep(20000000)
		}
	}()
}

func NewGUI(board *Board) *GUI {
	gui := new(GUI)
	X, _ := xgbutil.NewConn()
	gui.canvas = xgraphics.New(X, image.Rect(0, 0, board.Width(), board.Height()))
	gui.queue = make(map[Pos]bool, board.Width()*board.Height())
	gui.win = gui.canvas.XShow()
	gui.board = board
	board.OnUpdate = func(p Pos) { gui.Update(p) }
	gui.StartLoop()
	return gui
}
