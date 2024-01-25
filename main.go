package main

import (
	"tpt_backend/db"
	"tpt_backend/routes"
)

func main() {
	db.Init()

	e := routes.Init()

	e.Logger.Fatal(e.Start(":4000"))
}
