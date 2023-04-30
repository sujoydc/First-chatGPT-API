package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	var name string
	fmt.Print("What's your name? -")
	fmt.Scanln(&name)

	fmt.Println("You name recorded as: " + name)
}
