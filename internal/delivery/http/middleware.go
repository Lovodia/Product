package httpDelivery

import "net/http"

func errorHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			_ = renderError(w, r, err)
		}
	}
}
