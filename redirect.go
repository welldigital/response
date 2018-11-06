package response

import "net/http"

// Redirect sets up a redirect to the client.
func Redirect(url string, w http.ResponseWriter, status int) {
	w.Header().Set("Location", url)
	w.WriteHeader(status)
	w.Write([]byte{})
}
