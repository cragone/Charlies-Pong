package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

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

func main() {
	var cmd *exec.Cmd
	GameControls := CommandBoard{
		Cmd: cmd,
	}
	//Load into game
	fmt.Println("WELCOME TO PING PONG")
	time.Sleep(1 * time.Second)
	fmt.Println("3")
	time.Sleep(1 * time.Second)
	fmt.Println("2")
	time.Sleep(1 * time.Second)
	fmt.Println("1")
	time.Sleep(1 * time.Second)
	fmt.Println("GO!")

	//Declare the game
	PingPong := Game{
		GameDuration: 20,
		GameCommands: &GameControls,
		PlayerOne:    &Player{},
		PlayerTwo:    &Player{},
		GameBoard:    &Board{},
	}

	//Clear screen for game
	PingPong.ClearTerminal()
	// for user input

	//Build Game Board onto screen
	PingPong.CreateBoard(20, 120)
	PingPong.VolleyStart()

	PingPong.PlayerOne.Score = 0
	PingPong.PlayerTwo.Score = 0

	//run the Game timer in the background
	//we need a done chan to wait for the complete of the game timer
	//when the done chan receives a value it will exit the routines
	doneChan := make(chan bool, 1)
	PingPong.GameTimer(doneChan)
	inputChan := make(chan string)
	go ReadSingleKey(inputChan)

	//need a go routine which waits for key strokes in the background

	for {
		select {
		case done := <-doneChan:
			if done {
				log.Printf("Player One Scored: %d", PingPong.PlayerOne.Score)
				log.Printf("Player Two Scored: %d", PingPong.PlayerTwo.Score)
				log.Fatalf("GAME OVER!!!")
			} else {
				log.Fatalf("USER ENDED GAME")
			}
		case input := <-inputChan:
			PingPong.MovePlayer(input)
			PingPong.ScreenWriter()
			//used to write current standings and erase old ones
		}
	}

	//needs to await the doneChan's response

	///done is true the game timer ran to zero

}

// gametimer must run in the background of the game on a separate go routine
func (g *Game) GameTimer(doneChan chan bool) {
	go func() {
		time.Sleep(g.GameDuration * time.Second)
		g.ClearTerminal()
		doneChan <- true
	}()
}

// logic for hitting ball, or scoring points
func (g *Game) PlayerTwoHitBall() bool {
	return g.PlayerTwo.X == g.GameBall.X && g.PlayerTwo.Y == g.GameBall.Y
}

func (p *Player) GivePlayerPoint() {
	p.PlayerLock.Lock()
	defer p.PlayerLock.Unlock()
	p.Score += 1
}

func (g *Game) PlayerOneHitBall() bool {
	return g.PlayerOne.X == g.GameBall.X && g.PlayerOne.Y == g.GameBall.Y
}

func (g *Game) BallMovement(direction int) {
	go func() {
		slopeCalc := g.RandomLineGenerator()
		count := 0
		for {
			select {
			case <-g.BallStopChan:
				return
			default:
			}
			count++
			// Predict new position
			newX := g.GameBall.X + direction // if y = 1/3 every 3 x means 1 y up
			newY := g.GameBall.Y
			if count == slopeCalc {
				count = 0 //reset
				newY += 1 // move ball up one
			}

			// Check if someone scores
			if newX <= -1 {
				g.StopOnce.Do(func() { close(g.BallStopChan) })
				g.PlayerTwo.GivePlayerPoint()
				g.ResetBoard()
				g.VolleyStart()
				return
			}
			if newX >= g.GameBoard.Width {
				g.StopOnce.Do(func() { close(g.BallStopChan) })
				g.PlayerOne.GivePlayerPoint()
				g.ResetBoard()
				g.VolleyStart()
				return
			}

			// Check if hit by player
			if direction > 0 && g.PlayerTwoHitBall() {
				direction = -1
			}
			if direction < 0 && g.PlayerOneHitBall() {
				direction = 1
			}

			time.Sleep(10 * time.Millisecond)
			g.ScreenWriter()

			oldY, oldX := g.GameBall.Y, g.GameBall.X

			// Move ball
			g.GameBoard.BoardLock.Lock()
			spaceState := g.SaveOldSpaceState(newY, newX)
			g.GameBoard.Layout[newY][newX] = g.GameBoard.Layout[oldY][oldX]
			g.GameBoard.Layout[oldY][oldX] = spaceState
			g.GameBoard.BoardLock.Unlock()

			// Update position
			g.GameBall.BallLock.Lock()
			g.GameBall.Y, g.GameBall.X = newY, newX
			g.GameBall.BallLock.Unlock()
		}
	}()
}

// the problem you must solve is that y must always be less than the height
// and y must always be greater than 0
// the equation for a line is y = Mx + B
// based on current y it must not go outside that range at ending y
func (g *Game) RandomLineGenerator() int {
	//a number negative for up and a number positive for down and the range
	//must keep it within the boundaries
	spacesAboveLeft := g.GameBall.Y - g.GameBoard.Height //the out come must not move up more than this
	spacesBelowLeft := g.GameBoard.Height - g.GameBall.Y //the out come must not move below this
	//those are the ys in the slope equation
	//the returned number is how many spaces
	var possibleEndPoints []int
	// upDownDirector := rand.Intn(2)

	for endPoint := spacesAboveLeft; endPoint < spacesBelowLeft; endPoint++ {
		possibleEndPoints = append(possibleEndPoints, endPoint)
	}

	b := rand.Intn(len(possibleEndPoints))

	return b
} //y = mx+b

// func (g *Game) SpacesBelow() []int{}

func (g *Game) SaveOldSpaceState(y, x int) string {
	return g.GameBoard.Layout[y][x]
}

// writes the current game board to the screen
func (g *Game) ScreenWriter() {
	g.ClearTerminal()
	g.PrintCurrentGamePositions()
}

// this is the volley start method which sets the players and the ball
// it is called at the start of the game and can be called again to reset the game
func (g *Game) VolleyStart() {
	g.SetPlayerStart()
	g.SetPingPongBall()
	g.BallMovement(1)
	g.BallStopChan = make(chan struct{})
	g.StopOnce = sync.Once{}

}

// player methods
func (g *Game) SetPlayerStart() {
	g.GameBoard.BoardLock.Lock()
	g.PlayerOne.PlayerLock.Lock()
	g.PlayerTwo.PlayerLock.Lock()
	defer g.GameBoard.BoardLock.Unlock()
	defer g.PlayerOne.PlayerLock.Unlock()
	defer g.PlayerTwo.PlayerLock.Unlock()

	//left player
	g.PlayerOne.X = 1
	g.PlayerOne.Y = g.GameBoard.Height / 2
	//right player
	g.PlayerTwo.X = g.GameBoard.Width - 2
	g.PlayerTwo.Y = g.GameBoard.Height / 2
	//place the players onto the Board
	g.GameBoard.Layout[g.PlayerOne.Y][g.PlayerOne.X] = "X"
	g.GameBoard.Layout[g.PlayerTwo.Y][g.PlayerTwo.X] = "X"
	//return the players

}

// starting the ball next to player one maybe later we can change the logic to go to loser side
func (g *Game) SetPingPongBall() {
	gameBall := Ball{
		X: g.PlayerOne.X + 1, // place the ball next to player one
		Y: g.PlayerOne.Y,
	}
	g.GameBall = &gameBall                               // create a ball for the game
	g.GameBoard.Layout[g.GameBall.Y][g.GameBall.X] = "0" //place the ball on the board
}

// Move player
func (g *Game) MovePlayer(input string) {
	g.PlayerOne.PlayerLock.Lock()
	g.PlayerTwo.PlayerLock.Lock()
	defer g.PlayerOne.PlayerLock.Unlock()
	defer g.PlayerTwo.PlayerLock.Unlock()
	switch input {
	case "s":
		if g.PlayerOne.Y < g.GameBoard.Height-2 { //move player one down if not at border
			oldY := g.PlayerOne.Y
			newY := g.PlayerOne.Y + 1

			g.GameBoard.BoardLock.Lock()
			// move player marker
			g.GameBoard.Layout[newY][g.PlayerOne.X] = g.GameBoard.Layout[oldY][g.PlayerOne.X]
			// clear old position
			g.GameBoard.Layout[oldY][g.PlayerOne.X] = " "
			g.GameBoard.BoardLock.Unlock()

			// update player position
			g.PlayerOne.Y = newY
		}
	case "w":
		if g.PlayerOne.Y > 1 { //move player one up if not at the height
			oldY := g.PlayerOne.Y
			newY := g.PlayerOne.Y - 1

			g.GameBoard.BoardLock.Lock()
			g.GameBoard.Layout[newY][g.PlayerOne.X] = g.GameBoard.Layout[oldY][g.PlayerOne.X]
			g.GameBoard.Layout[oldY][g.PlayerOne.X] = " "
			g.GameBoard.BoardLock.Unlock()

			// update the player position
			g.PlayerOne.Y = newY

		}
	case "l":
		if g.PlayerTwo.Y < g.GameBoard.Height-2 { //move player two down if not at border
			oldY := g.PlayerTwo.Y
			newY := g.PlayerTwo.Y + 1

			g.GameBoard.BoardLock.Lock()
			// move player marker
			g.GameBoard.Layout[newY][g.PlayerTwo.X] = g.GameBoard.Layout[oldY][g.PlayerTwo.X]
			// clear old position
			g.GameBoard.Layout[oldY][g.PlayerTwo.X] = " "
			g.GameBoard.BoardLock.Unlock()

			// update player position
			g.PlayerTwo.Y = newY
		}

	case "o":
		if g.PlayerTwo.Y > 1 {
			oldY := g.PlayerTwo.Y
			newY := g.PlayerTwo.Y - 1

			g.GameBoard.BoardLock.Lock()
			//move player marker
			g.GameBoard.Layout[newY][g.PlayerTwo.X] = g.GameBoard.Layout[oldY][g.PlayerTwo.X]
			//clear the old position
			g.GameBoard.Layout[oldY][g.PlayerTwo.X] = " "
			g.GameBoard.BoardLock.Unlock()

			//update player position
			g.PlayerTwo.Y = newY
		}
	}
}

func (g *Game) ResetBoard() {
	g.ClearTerminal()
	g.CreateBoard(g.GameBoard.Height, g.GameBoard.Width)
}

// board methods
// Build Game Board Builds out the game in strings using height and width
func (g *Game) CreateBoard(height, width int) {
	if g.GameBoard.Layout == nil || len(g.GameBoard.Layout) != height {
		g.GameBoard.Layout = make([][]string, height)
	}

	for y := 0; y < height; y++ {
		if g.GameBoard.Layout[y] == nil || len(g.GameBoard.Layout[y]) != width {
			g.GameBoard.Layout[y] = make([]string, width)
		}

		for x := 0; x < width; x++ {
			if x == 0 || x == width-1 {
				g.GameBoard.Layout[y][x] = "|"
			} else if y == 0 || y == height-1 {
				g.GameBoard.Layout[y][x] = "-"
			} else {
				g.GameBoard.Layout[y][x] = " "
			}
		}
	}

	// Update dimensions (optional if they might change)
	g.GameBoard.Width = width
	g.GameBoard.Height = height
}

// Prints the current board
func (g *Game) PrintCurrentGamePositions() {
	for _, row := range g.GameBoard.Layout {
		fmt.Println(strings.Join(row, ""))
	}
}

// over arching game methods

// ClearTerminal to clear existing UI
func (g *Game) ClearTerminal() {
	g.GameCommands.Mu.Lock()
	defer g.GameCommands.Mu.Unlock()
	g.GameCommands.Cmd = exec.Command("cmd", "/c", "cls")
	g.GameCommands.Cmd.Stdout = os.Stdout
	g.GameCommands.Cmd.Run()
}

// used to read key strokes sits in a go routine and reads to a channel
func ReadSingleKey(inputChan chan string) {
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b []byte = make([]byte, 1)
	for {
		if _, err := os.Stdin.Read(b); err == nil {
			inputChan <- string(b[0])
		}
	}
}
