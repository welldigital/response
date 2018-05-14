package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes the value v as JSON to the ResponseWriter.
func JSON(v interface{}, w http.ResponseWriter, status int) (err error) {
	data, err := json.Marshal(v)
	if err != nil {
		ErrorString("json_marshal", "failed to marshal JSON", w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
	return
}

// List is the return type for JSON lists.
type List struct {
	StartIndex int         `json:"startIndex"`
	PageSize   int         `json:"pageSize"`
	NextIndex  *int        `json:"nextIndex,omitempty"`
	Data       interface{} `json:"data"`
}

// NewList creates a JSON list type.
func NewList(startIndex, pageSize, nextIndex int, data interface{}) List {
	l := List{
		StartIndex: startIndex,
		PageSize:   pageSize,
		Data:       data,
	}
	if nextIndex > 0 {
		l.NextIndex = &nextIndex
	}
	return l
}
