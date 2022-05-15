package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const defaultPort = 8080

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/slow", handleRequest)

	log.Println("Listening on port:", defaultPort)

	err := http.ListenAndServe(fmt.Sprintf(":%v", defaultPort), router)

	if err != nil {
		log.Println("Error starting server")
	}

}

func handleRequest(w http.ResponseWriter, req *http.Request) {

	s := getSlowData()

	b := strings.Join(s, "\n")

	w.Write([]byte(b))
}

func getSlowData() []string {
	// artificial delay
	time.Sleep(time.Second * 2)

	num := 100
	res := make([]string, num)
	for i := 0; i < num; i++ {
		res[i] = fmt.Sprintf("dummy data line %v", i)
	}

	return res
}
