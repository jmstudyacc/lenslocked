// tells Go - package main is the entrypoint and main() is where the code starts
package main

// importing fmt to print out stuff and net/http to set up webserver
import (
	"fmt"
	"net/http"
)

// creating a function to handle the http request & present a web page
func handlerFunc(writer http.ResponseWriter, reader *http.Request) {
	fmt.Fprint(writer, "<h1>Welcome to the awesome site!</h1>")
}

func main() {
	// registration - pass in a pattern to determine the paths to be used by the handler and a function for the handler
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the webserver on :3000")
	// starts the server & prevents the code from exiting
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		panic(err)
	}
}
