package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	notesPath       = "/v1alpha1/projects/image-signing/notes"
	occurrencesPath = "/v1alpha1/projects/image-signing/occurrences"
)

func main() {
	// use PORT environment variable, or default to 8080
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	// register hello function to handle all requests
	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println(string(data))

		ar := v1beta1.AdmissionReview{}
		if err := json.Unmarshal(data, &ar); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pod := v1.Pod{}
		if err := json.Unmarshal(ar.Request.Object.Raw, &pod); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, container := range pod.Spec.Containers {
			log.Println("container Image is", container.Image)

			if !strings.Contains(container.Image,"docker.jfrog.io") {
				log.Println("Image is not being pulled from Private Registry")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}


		log.Printf("Serving request: %s", r.URL.Path)
	})

	server.HandleFunc("/ping", hello)


	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

// ping responds to the request with a plain-text "Ok" message.
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	fmt.Fprintf(w, "Ok")
}