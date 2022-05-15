# go-progressive-download-poc

```go
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"vidhill.com/dummy/streamingserver/model"
)

const port = 8080

func main() {

	rand.Seed(time.Now().UnixNano())

	// pr, pw := io.Pipe()

	// go ThisIsJustARoutineProcedureDontWorry(pw)

	// io.Copy(os.Stdout, pr)

	mux := chi.NewMux()

	mux.Get("/slow-csv", handleCSV)

	log.Println("Listening on port:", port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)

	if err != nil {
		log.Println("Error starting server")
	}

}

func handleCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename=assetChecklistData.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Transfer-Encoding", "chunked")

	pr, pw := io.Pipe()

	writer := ConvertHttpWriterToGzipWriter(w)

	defer writer.Close()

	go ThisIsARoutineProcedureSirDontWorry(pw)

	io.Copy(writer, pr)
	// io.Copy(os.Stdout, pr)

}

// func ThisIsARoutineProcedureSirDontWorry(pw *io.PipeWriter, flusher http.Flusher) {
func ThisIsARoutineProcedureSirDontWorry(w *io.PipeWriter) {

	// close the writer, so the reader knows that there's no more data
	defer w.Close()

	b := new(bytes.Buffer)
	csvWriter := csv.NewWriter(b)

	numInvocations := makeRandRange(100, 150)
	log.Println("numInvocations", numInvocations)

	for i := 0; i < numInvocations; i++ {

		dummyCSV := model.MakeFakeData(i, 1_000)

		csvWriter.WriteAll(dummyCSV)

		// w.Write(b.Bytes())

		io.Copy(w, b)

		b.Reset() // clear any data the bufffer

		// artificial delay
		time.Sleep(time.Second * 3)
	}

}

func ConvertHttpWriterToGzipWriter(w http.ResponseWriter) *gzip.Writer {
	writer, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return writer
}

func makeRandRange(min, max int) int {
	return rand.Intn(max-min) + min
}


```
