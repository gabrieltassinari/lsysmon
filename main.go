package main

import (
	"fmt"
	"net/http"

	"github.com/gabrieltassinari/lsysmon/routes"
)

func main() {
	routes.Routes()

	fmt.Println("Running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
