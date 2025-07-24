package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

type Methods interface {
	SetPlayerPosition(board *Board)
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
		GameDuration: 10,
		GameCommands: &GameControls,
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

	//need a go routine which waits for key strokes in the background
	inputChan := make(chan string)
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
		default:
			go func() {
				ReadSingleKey(inputChan) //running in the background to receiving inputs
			}()
			// PingPong.BallStopChan = make(chan struct{})
			PingPong.BallMovement(1)
			input := <-inputChan
			PingPong.MovePlayer(input)
			PingPong.ScreenWriter() //used to write current standings and erase old ones
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

func (g *Game) PlayerTwoScores() bool {
	return g.GameBall.X == g.GameBoard.Width-2
}

func (p *Player) GivePlayerPoint() {
	p.Score += 1
}

func (g *Game) PlayerOneScores() bool {
	return g.GameBall.X == 1
}

func (g *Game) PlayerOneHitBall() bool {
	return g.PlayerOne.X == g.GameBall.X && g.PlayerOne.Y == g.GameBall.Y
}

func (g *Game) BallMovement(direction int) {
	go func() {
		g.GameBall.BallLock.Lock()
		defer g.GameBall.BallLock.Unlock()

		for {
			select {
			case <-g.BallStopChan:
				return
			default:
			}

			// Predict new position
			newX := g.GameBall.X + direction
			newY := g.GameBall.Y

			// Check if someone scores
			if newX <= 0 {
				g.StopOnce.Do(func() {
					close(g.BallStopChan)
				})
				g.PlayerTwo.GivePlayerPoint()
				g.ResetBoard()
				g.VolleyStart()
				return
			}
			if newX >= g.GameBoard.Width-1 {
				g.StopOnce.Do(func() {
					close(g.BallStopChan)
				})
				g.PlayerOne.GivePlayerPoint()
				g.ResetBoard()
				g.VolleyStart()
				return
			}

			// Check if hit by player
			if direction > 0 && g.PlayerTwoHitBall() {
				g.BallMovement(-1)
				return
			}
			if direction < 0 && g.PlayerOneHitBall() {
				g.BallMovement(1)
				return
			}

			time.Sleep(10 * time.Millisecond)
			g.ScreenWriter()

			oldY, oldX := g.GameBall.Y, g.GameBall.X

			// Move ball
			g.GameBoard.BoardLock.Lock()
			g.GameBoard.Layout[newY][newX] = g.GameBoard.Layout[oldY][oldX]
			g.GameBoard.Layout[oldY][oldX] = " "
			g.GameBoard.BoardLock.Unlock()

			// Update position
			g.GameBall.Y, g.GameBall.X = newY, newX
		}
	}()
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
	g.PrintCurrentGamePositions()
	g.BallStopChan = make(chan struct{})
	g.StopOnce = sync.Once{}

}

// player methods
func (g *Game) SetPlayerStart() {
	g.GameBoard.BoardLock.Lock()
	defer g.GameBoard.BoardLock.Unlock()
	//left player
	player1 := Player{
		X: 1, //needs to stay within the border at 0
		Y: g.GameBoard.Height / 2,
	}
	g.PlayerOne = &player1
	//right player
	player2 := Player{
		X: g.GameBoard.Width - 2, //needs to stay within the border at width -1
		Y: g.GameBoard.Height / 2,
	}
	g.PlayerTwo = &player2
	//place the players onto the Board
	g.GameBoard.Layout[player1.Y][player1.X] = "X"
	g.GameBoard.Layout[player2.Y][player2.X] = "X"
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
	g.GameBoard.Layout = [][]string{}
	g.ClearTerminal()
	g.CreateBoard(g.GameBoard.Height, g.GameBoard.Width)
}

// board methods
// Build Game Board Builds out the game in strings using height and width
func (g *Game) CreateBoard(height, width int) {
	layout := make([][]string, height)
	for y := 0; y < height; y++ {
		row := make([]string, width)
		for x := 0; x < width; x++ {
			if x == 0 { // left border
				row[x] = "|"
			} else if x == width-1 { // right border
				row[x] = "|"
			} else if y == 0 {
				row[x] = "-"
			} else if y == height-1 {
				row[x] = "-"
			} else {
				row[x] = " " //single game space
			}
		}
		layout[y] = row
	}
	g.GameBoard = &Board{
		Width:  width,
		Height: height,
		Layout: layout,
	}
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
	os.Stdin.Read(b)
	inputChan <- string(b[0])
}
