package graphics

import (
	"os/exec"
	"sync"
	"time"
)

func NewGame(d time.Duration, cmd *exec.Cmd) *Game {
	return &Game{
		GameDuration: 20,
		GameCommands: &CommandBoard{Cmd: cmd},
		PlayerOne:    &Player{},
		PlayerTwo:    &Player{},
		GameBoard:    &Board{},
	}
}

// players may need to be in the board and not the game...
type Game struct {
	GameBoard    *Board
	PlayerOne    *Player
	PlayerTwo    *Player
	GameBall     *Ball
	GameCommands *CommandBoard
	GameDuration time.Duration
	BallStopChan chan struct{}
	StopOnce     sync.Once
}

type Player struct {
	X          int
	Y          int
	Score      int
	PlayerLock sync.Mutex
}

type Ball struct {
	X        int
	Y        int
	BallLock sync.Mutex
}

type Board struct {
	Width     int
	Height    int
	Layout    [][]string
	BoardLock sync.Mutex
}

type CommandBoard struct {
	Cmd *exec.Cmd
	Mu  sync.Mutex
}
