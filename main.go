package main

import (
	"fmt"
	"forum/database"
	web "forum/web"
)

func main() {
	db := database.ConnectToDatabase()
	web.ConnectWeb(db)
	err := database.CreateTables(db)
	if err != nil {
		fmt.Println(err)
	}
}
