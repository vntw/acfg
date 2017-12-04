package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ClientAppHandler(handler http.Handler, httpfs http.FileSystem) http.HandlerFunc {
	idx := readIndexHTML(httpfs)

	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path
		f, err := httpfs.Open(filename)

		if err == nil {
			//log.Println("Access file:", filename)
			f.Close()
			handler.ServeHTTP(w, r)
			return
		}

		// Let the client app do the routing
		if os.IsNotExist(err) {
			//log.Println("Fall back to index.html:", filename)
			w.Write(idx)
			return
		}

		// Something is wrong, err != nil
		log.Printf("Could not get file %s, got error %v\n", filename, err)
		w.WriteHeader(500)
	}
}

func readIndexHTML(httpfs http.FileSystem) []byte {
	index, err := httpfs.Open("index.html")
	if err != nil {
		log.Fatal("Could not find index.html")
	}

	defer index.Close()

	idx, err := ioutil.ReadAll(index)
	if err != nil {
		log.Fatal("Could not read index.html")
	}

	return idx
}
