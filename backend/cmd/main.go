package main

import (
	"log"

	api "gitlab.atp-fivt.org/fullstack2024a/kondrashovti-project/backend/internal"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	api := api.NewApi("postgres://kondrashovti:admin@localhost:5432/programming_educ")

	api.UseCors("http://130.193.44.85")
	api.Register()
	api.Start(":5000")
	defer api.Close()
}
