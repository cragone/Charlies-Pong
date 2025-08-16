package main

import (
	"game/graphics"
	"game/tui"
	"log"
	"os/exec"
	"time"
)

func main() {
	var cmd *exec.Cmd
	//Load into game
	tui.GameIntro()

	pingPong := graphics.NewGame(60*time.Second, cmd)

	//Declare the game

	//Clear screen for game
	pingPong.ClearTerminal()
	// for user input

	//Build Game Board onto screen
	pingPong.CreateBoard(20, 120)
	pingPong.VolleyStart()

	pingPong.PlayerOne.Score = 0
	pingPong.PlayerTwo.Score = 0

	//run the Game timer in the background
	//we need a done chan to wait for the complete of the game timer
	//when the done chan receives a value it will exit the routines
	doneChan := make(chan bool, 1)
	pingPong.GameTimer(doneChan)
	inputChan := make(chan string)
	go tui.ReadSingleKey(inputChan)

	//need a go routine which waits for key strokes in the background
	// go func() {
	// 	mux := http.NewServeMux()

	// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{"message": "hello"})
	// 	})

	// 	http.ListenAndServe(":80", mux)
	// }()

	for {
		select {
		case done := <-doneChan:
			if done {
				log.Printf("Player One Scored: %d", pingPong.PlayerOne.Score)
				log.Printf("Player Two Scored: %d", pingPong.PlayerTwo.Score)
				log.Fatalf("GAME OVER!!!")
			} else {
				log.Fatalf("USER ENDED GAME")
			}
		case input := <-inputChan:
			pingPong.MovePlayer(input)
			pingPong.ScreenWriter()
			//used to write current standings and erase old ones
		}
	}

	//needs to await the doneChan's response

	///done is true the game timer ran to zero

}
