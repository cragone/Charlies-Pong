package tui

import (
	"fmt"
	"time"
)

func GameIntro() {
	fmt.Println("WELCOME TO PING PONG")
	time.Sleep(1 * time.Second)
	fmt.Println("3")
	time.Sleep(1 * time.Second)
	fmt.Println("2")
	time.Sleep(1 * time.Second)
	fmt.Println("1")
	time.Sleep(1 * time.Second)
	fmt.Println("GO!")
}
