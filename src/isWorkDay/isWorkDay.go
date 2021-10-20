package isWorkDay

import (
	"fmt"
	"net/http"
)

func WorkDay(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "<h1>Hello from Go!</h1>")
}
