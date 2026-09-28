package main

import (
	"fmt"
)

var choice string

func main() {
	fmt.Println("Welcome to the Docker TUI: ")

	for {
		tool := DockerTool{Name: "Docker TUI"}

		fmt.Print(`
Please enter a valid menu option:

1: Check Docker Containers
2: Change Existing Docker Container State
3: Create Docker Container
Q: Quit

>  `)
		fmt.Scan(&choice)

		switch choice {

		case "1":
			fmt.Println("Checking Docker Containers...")
			tool.printContainerList()
			continue

		case "2":
			fmt.Println("Changing Existing Docker Container State...")
			tool.changeExistingContainerState()
			continue

		case "3":
			fmt.Println("Creating New Docker Container...")
			tool.createNewContainer()

		case "q", "Q":
			fmt.Println("Quitting...")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
