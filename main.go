package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/whoismath/ssrf-pwnabot/app"
	"github.com/whoismath/ssrf-pwnabot/config"
)

func secretHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := config.Setup()

	fmt.Fprintf(w, "Hacked by Ganesh<img src=\"%s\"></img>", c.Secret)
}

func main() {

	c, err := config.Setup()
	if err != nil {
		log.Fatal(err)
	}

	// Secret service
	http.HandleFunc("/", secretHandler)
	go http.ListenAndServe(":8080", nil)

	a := app.CreateApp(c)
	a.Start()

}
