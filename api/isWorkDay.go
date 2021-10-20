package handler

import (
	"fmt"
	"net/http"
)

func Handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "<h1>Hello from Go!</h1>")
}
