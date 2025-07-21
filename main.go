package main

import (
	"log"
)

func main() {
	app := App{}
	log.Println("Server starting on :1481")
	app.Initialize()
	app.Run("localhost:1481")
}
