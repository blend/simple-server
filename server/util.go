package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, obj interface{}) {
	contents, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(contents))
}

func writeBinary(w http.ResponseWriter, content []byte) {
	fmt.Fprintln(w, base64.RawURLEncoding.EncodeToString(content))
}

func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
