package main

import (
	"fmt"
	"os"

	"github.com/KiranPawar0/Data-Drive-System/pkg/database"
	"github.com/KiranPawar0/Data-Drive-System/routes/task"
)

func main() {
	dbConn, err := database.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)

		return
	}

	var serviceName string

	// Check if the SERVICE_NAME environment variable is set
	if envServiceName := os.Getenv("SERVICE_NAME"); envServiceName != "" {
		serviceName = envServiceName
	} else {
		// Check if a command-line argument is provided
		if len(os.Args) < 2 {
			fmt.Println("Service name not provided")
			return
		}
		serviceName = os.Args[1]
	}

	switch serviceName {
	case "task":
		task.Task(dbConn) // Start the auth server
	}
}
