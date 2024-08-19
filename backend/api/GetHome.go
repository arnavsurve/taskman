package api

import (
	"fmt"
	"io"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /")
	io.WriteString(w, "Home\n")
}
