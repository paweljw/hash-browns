package main

import "crypto/md5"
import "net/http"
import "log"
import "fmt"
import "io"
import "encoding/hex"
import "os"

func urlMd5(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, response.Body); err != nil {
		log.Fatal(err)
	}
	sum := hex.EncodeToString(hash.Sum(nil))

	return sum, nil
}

func sumResponse(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["q"]

		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(422)
			log.Println("Url Param 'q' is missing")
			fmt.Fprintf(w, "Url Param 'q' is missing")
			return
		}

		key := keys[0]
		md5, err := urlMd5(key)

		if err != nil {
			w.WriteHeader(422)
			fmt.Fprintf(w, "%s", err)
		} else {
			fmt.Fprintf(w, "%s", md5)
		}
}

func main () {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/sum", sumResponse)

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
