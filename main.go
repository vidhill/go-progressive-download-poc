package main

import (
	"bytes"
	"fmt"
	"io"
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
	// indicate to the client to expect a chunked response
	w.Header().Set("Transfer-Encoding", "chunked")

	pr, pw := io.Pipe()

	go ThisIsARoutineProcedureSirDontWorry(pw)

	io.Copy(w, pr)

	log.Println("slow request completed")
}

func ThisIsARoutineProcedureSirDontWorry(w io.WriteCloser) {

	b := new(bytes.Buffer)

	numIterations := 10

	for i := 0; i < numIterations; i++ {
		// simulate a slow operation that
		// for some hypothetical reason can only be done sequentially
		s := getSlowData(i)

		lines := strings.Join(s, "\n")

		// write the partial data to the Writer
		// as it is retrieved
		w.Write([]byte(lines + "\n"))

		// clear any data the buffer
		b.Reset()

		log.Printf("iteration %v/%v complete", i+1, numIterations)

	}

	// when all sequential calls are done
	// close the writer, so the reader knows not to expect any more data
	b.Reset()
	w.Close()
}

//
// fake slow data
//
func getSlowData(iteration int) []string {
	// artificial delay
	time.Sleep(time.Second * 2)

	num := 1000
	res := make([]string, num)

	for i := 0; i < num; i++ {
		res[i] = fmt.Sprintf("dummy data line, iteration %v line %v", iteration, i)
	}

	return res
}
