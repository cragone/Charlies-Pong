package graphics

import "time"

// gametimer must run in the background of the game on a separate go routine
func (g *Game) GameTimer(doneChan chan bool) {
	go func() {
		time.Sleep(g.GameDuration * time.Second)
		g.ClearTerminal()
		doneChan <- true
	}()
}
