package response

import (
	"encoding/xml"
	"net/http"
)

// XML writes the value v as XML to the ResponseWriter.
func XML(v interface{}, w http.ResponseWriter, status int) (err error) {
	data, err := xml.Marshal(v)
	if err != nil {
		ErrorString("xml_marshal", "failed to marshal XML", w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	w.Write(data)
	return
}
