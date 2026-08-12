package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isPOSTAndJSON(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprint(w, "Неправильный метод")
		return false
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(415)
		fmt.Fprint(w, "Unsupported Media Type")
		return false
	}
	return true
}

func mapErrors(err error, w http.ResponseWriter) {
	st, ok := status.FromError(err)
	if ok && st.Code() == codes.InvalidArgument {
		w.WriteHeader(400)
		fmt.Fprint(w, st.Message())
		return
	}
	if ok && st.Code() == codes.AlreadyExists {
		w.WriteHeader(409)
		fmt.Fprint(w, st.Message())
		return
	}
	if ok && st.Code() == codes.Unauthenticated {
		w.WriteHeader(401)
		fmt.Fprint(w, st.Message())
		return
	}
	log.Printf("error: %v", err)
	w.WriteHeader(500)
	fmt.Fprint(w, "internal error")
}

func writeJSON(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
