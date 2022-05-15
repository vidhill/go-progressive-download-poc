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
}

func ThisIsARoutineProcedureSirDontWorry(w *io.PipeWriter) {

	b := new(bytes.Buffer)

	for i := 0; i < 10; i++ {

		s := getSlowData(i)

		lines := strings.Join(s, "\n")

		// write the partial data to the Pipewriter
		// as it is retrieved
		w.Write([]byte(lines + "\n"))

		// clear any data the bufffer
		b.Reset()

	}

	// when all sequential calls are done
	// close the writer, so the reader knows not to expect any more data
	w.Close()

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
