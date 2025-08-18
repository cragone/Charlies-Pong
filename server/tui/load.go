package tui

import (
	"fmt"
	"time"
)

// Game Start Count Down Messages
const (
	WelcomingMessage string = "WELCOME TO PING PONG"
	GameCountDown3   string = "3"
	GameCountDown2   string = "2"
	GameCountDown1   string = "1"
	BeginingMessage  string = "GO!!!"
)

func GameIntro() {
	fmt.Println(WelcomingMessage)
	time.Sleep(1 * time.Second)
	fmt.Println(GameCountDown3)
	time.Sleep(1 * time.Second)
	fmt.Println(GameCountDown2)
	time.Sleep(1 * time.Second)
	fmt.Println(GameCountDown1)
	time.Sleep(1 * time.Second)
	fmt.Println(BeginingMessage)
}
