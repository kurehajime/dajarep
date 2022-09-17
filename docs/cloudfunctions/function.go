// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kurehajime/dajarep"
)

func Dajareper(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Messages []string `json:"messages"`
	}
	var allows = []string{"http://127.0.0.1:8080", "https://kurehajime.github.io"}
	for _, v := range allows {
		if v == r.Header.Get("Origin") {
			w.Header().Set("Access-Control-Allow-Origin", v)
			break
		}
	}

	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	input := strings.Join(d.Messages, "\n")
	dajares, _ := dajarep.Dajarep(input)
	d.Messages = dajares
	out, err := json.Marshal(d)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json,;charset=UTF-8")
	w.Write(out)
}
