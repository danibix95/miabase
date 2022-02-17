package response

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSON is a convenient function to write a JSON object as response to the incoming request
func JSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, `{"message":"error encoding response"}`); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
