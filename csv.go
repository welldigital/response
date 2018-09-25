package response

import (
	"encoding/csv"
	"net/http"
)

// CSV writes v as a CSV file to the ResponseWriter.``
func CSV(v [][]string, w http.ResponseWriter, status int) (err error) {
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(status)
	csvw := csv.NewWriter(w)
	csvw.WriteAll(v)
	csvw.Flush()
	return
}
