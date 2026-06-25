package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome Veerbal to golang backend journey")
	})

	fmt.Println("Server is listening on port :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
