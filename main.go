package main

import (
	"fmt"
	"log"
	"net/http"

	router "github.com/mfc-creations/mongoapi/routes"
)

func main() {
	fmt.Println("Mongodb Api")
	r:=router.Router()
	fmt.Println("Server is getting started");
	log.Fatal(http.ListenAndServe(":4000",r))
	fmt.Println("Listening at port 4000")
}
