package main

import "fmt"

// 程序的入口点
func main() {
	fmt.Println("Hello world!")
	say("Hello Go!")
}
func say(message string) {
	fmt.Println("You said: ", message)
}
