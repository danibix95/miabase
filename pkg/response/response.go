package response

import (
	"encoding/json"
	"io"
	"net/http"
)

func JSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, `"msg":"error encoding response"`); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
