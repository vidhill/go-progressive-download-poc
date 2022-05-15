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

// handler
func handleRequest(w http.ResponseWriter, req *http.Request) {

	agg := []string{}

	for i := 0; i < 10; i++ {
		// slow running fetch,
		// which for whatever reason must be retrieved sequentially
		s := getSlowData(i)

		agg = append(agg, s...)
	}

	lines := strings.Join(agg, "\n")

	w.Write([]byte(lines))
}

//
// fake slow data
//
func getSlowData(iteration int) []string {
	// artificially delay
	time.Sleep(time.Second * 2)

	num := 100
	res := make([]string, num)

	for i := 0; i < num; i++ {
		res[i] = fmt.Sprintf("dummy data line, iteration %v line %v", iteration, i)
	}

	return res
}
