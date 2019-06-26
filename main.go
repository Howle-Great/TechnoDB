package main

import (
	"./dbhandlers"
	"./handlers"
	"fmt"
)

func main() {
	dbhandlers.DB.Connetc()
	router := handlers.CreateRouter()	
	router.Run(":5000")
	fmt.Println("Starting server on 127.0.0.1:5000")
}