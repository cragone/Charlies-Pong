package graphics

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

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

func BallMovementOnYAxis(randomEndingY int) int {
	if randomEndingY < 0 {
		return Up
	}
	return Down
}

func (g *Game) BallMovement(direction int) {
	go func() {
		// N = how many X steps before we move Y by dy (±1)
		N := g.RandomEndingYGenerator()
		if N == 0 {
			N = 1
		}
		if N < 0 {
			N = -N
		}

		dy := BallMovementOnYAxis(N) // must be either +1 or -1
		steps := 0
		dir := direction // local copy so we can flip safely

		for {
			select {
			case <-g.BallStopChan:
				return
			default:
			}

			steps++

			// Predict new position
			newX := g.GameBall.X + dir
			newY := g.GameBall.Y

			// Move vertically every N x-steps
			if steps >= N {
				steps = 0
				newY += dy
			}

			// Reflect off top/bottom and keep Y in-bounds
			boardH := g.GameBoard.Height
			if newY < 0 {
				newY = 0
				dy = +1 // bounce downward
			} else if newY >= boardH {
				newY = boardH - 1
				dy = -1 // bounce upward
			}

			// Check if someone scores (left/right walls)
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

			// Paddle collisions flip horizontal direction.
			// (If you also want to change slope on paddle hit, regenerate N/dy here)
			if dir == Right && g.PlayerTwoHitBall() {
				dir = Left
				// Optional: change vertical cadence/slope on hit:
				// N = max(1, abs(g.RandomEndingYGenerator()))
				// dy = BallMovementOnYAxis(N)
			}
			if dir == Left && g.PlayerOneHitBall() {
				dir = Right
				// Optional: change slope here too (see above).
			}

			time.Sleep(10 * time.Millisecond)
			g.ScreenWriter()

			// Commit move (guard indices!)
			g.GameBoard.BoardLock.Lock()
			oldY, oldX := g.GameBall.Y, g.GameBall.X

			// Save/restore space state
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

func RandomlyChooseUpOrDown() int {
	if rand.Intn(2) == 0 {
		return Up
	}
	return Down

}

// the problem you must solve is that y must always be less than the height
// and y must always be greater than 0
// the equation for a line is y = Mx + B
// based on current y it must not go outside that range at ending y
func (g *Game) RandomEndingYGenerator() int {
	//a number negative for up and a number positive for down and the range
	//must keep it within the boundaries
	h := g.GameBoard.Height
	y := g.GameBall.Y

	if h <= 0 {
		return 0
	}
	// Max steps available without crossing bounds
	maxUp := y           // can move up at most 'y' steps to hit 0
	maxDown := h - 1 - y // can move down at most 'h-1-y' steps to hit Height-1
	if maxUp < 0 {
		maxUp = 0
	}
	if maxDown < 0 {
		maxDown = 0
	}

	// Now choose a direction that has room; fall back if one side is blocked
	dir := RandomlyChooseUpOrDown()
	if (dir == Up && maxUp == 0) && maxDown > 0 {
		dir = Down
	} else if (dir == Down && maxDown == 0) && maxUp > 0 {
		dir = Up
	}
	// Produce a non-zero delta that stays in-bounds
	switch dir {
	case Up:
		if maxUp == 0 {
			return 0
		}
		// rand.Intn(n) yields [0, n), so +1 gives [1, n]
		return -(rand.Intn(maxUp) + 1)
	default: // Down
		if maxDown == 0 {
			return 0
		}
		return rand.Intn(maxDown) + 1
	} //to get that endpoint you need to move up or down
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
	g.BallMovement(Right)
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
	g.GameBoard.Layout[g.PlayerOne.Y][g.PlayerOne.X] = PlayerSymbol
	g.GameBoard.Layout[g.PlayerTwo.Y][g.PlayerTwo.X] = PlayerSymbol
	//return the players

}

// starting the ball next to player one maybe later we can change the logic to go to loser side
func (g *Game) SetPingPongBall() {
	gameBall := Ball{
		X: g.PlayerOne.X + Right, // place the ball next to player one
		Y: g.PlayerOne.Y,
	}
	g.GameBall = &gameBall                                      // create a ball for the game
	g.GameBoard.Layout[g.GameBall.Y][g.GameBall.X] = BallSymbol //place the ball on the board
}

// Move player logic

func (g *Game) MovePlayer(input string) {
	g.PlayerOne.PlayerLock.Lock()
	g.PlayerTwo.PlayerLock.Lock()
	defer g.PlayerOne.PlayerLock.Unlock()
	defer g.PlayerTwo.PlayerLock.Unlock()
	switch input {
	case PlayerOneDownKey:
		g.MovePlayerOneDown()
	case PlayerOneUpKey:
		g.MovePlayerOneUp()
	case PlayerTwoDownKey:
		g.MovePlayerTwoDown()
	case PlayerTwoUpKey:
		g.MovePlayerTwoUp()
	}
}

func (g *Game) MovePlayerOneUp() {
	if g.PlayerOne.Y > 1 { //move player one up if not at the height
		oldY := g.PlayerOne.Y
		newY := g.PlayerOne.Y + Up

		g.GameBoard.BoardLock.Lock()
		g.GameBoard.Layout[newY][g.PlayerOne.X] = g.GameBoard.Layout[oldY][g.PlayerOne.X]
		g.GameBoard.Layout[oldY][g.PlayerOne.X] = EmptySpace
		g.GameBoard.BoardLock.Unlock()

		// update the player position
		g.PlayerOne.Y = newY

	}
}

func (g *Game) MovePlayerOneDown() {
	if g.PlayerOne.Y < g.GameBoard.Height-2 { //move player one down if not at border
		oldY := g.PlayerOne.Y
		newY := g.PlayerOne.Y + Down

		g.GameBoard.BoardLock.Lock()
		// move player marker
		g.GameBoard.Layout[newY][g.PlayerOne.X] = g.GameBoard.Layout[oldY][g.PlayerOne.X]
		// clear old position
		g.GameBoard.Layout[oldY][g.PlayerOne.X] = EmptySpace
		g.GameBoard.BoardLock.Unlock()

		// update player position
		g.PlayerOne.Y = newY
	}
}

func (g *Game) MovePlayerTwoUp() {
	if g.PlayerTwo.Y > 1 {
		oldY := g.PlayerTwo.Y
		newY := g.PlayerTwo.Y + Up

		g.GameBoard.BoardLock.Lock()
		//move player marker
		g.GameBoard.Layout[newY][g.PlayerTwo.X] = g.GameBoard.Layout[oldY][g.PlayerTwo.X]
		//clear the old position
		g.GameBoard.Layout[oldY][g.PlayerTwo.X] = EmptySpace
		g.GameBoard.BoardLock.Unlock()

		//update player position
		g.PlayerTwo.Y = newY
	}
}

func (g *Game) MovePlayerTwoDown() {
	if g.PlayerTwo.Y < g.GameBoard.Height-2 { //move player two down if not at border
		oldY := g.PlayerTwo.Y
		newY := g.PlayerTwo.Y + Down

		g.GameBoard.BoardLock.Lock()
		// move player marker
		g.GameBoard.Layout[newY][g.PlayerTwo.X] = g.GameBoard.Layout[oldY][g.PlayerTwo.X]
		// clear old position
		g.GameBoard.Layout[oldY][g.PlayerTwo.X] = EmptySpace
		g.GameBoard.BoardLock.Unlock()

		// update player position
		g.PlayerTwo.Y = newY
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
				g.GameBoard.Layout[y][x] = VerticalBorder
			} else if y == 0 || y == height-1 {
				g.GameBoard.Layout[y][x] = HorizontalBorder
			} else {
				g.GameBoard.Layout[y][x] = EmptySpace
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
