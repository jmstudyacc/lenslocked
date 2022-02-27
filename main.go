// tells Go - package main is the entrypoint and main() is where the code starts
package main

// importing fmt to print out stuff and net/http to set up webserver
import (
	"fmt"
	"net/http"
)

// creating a function to handle the http request & present a web page
func homeHandler(writer http.ResponseWriter, reader *http.Request) {
	// defines the Content-Type header to be used >> Content Type Headers - HTTP Headers ~ Mozilla
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(writer, "<h1>Welcome to the awesome site!</h1>")
}

func contactHandler(writer http.ResponseWriter, reader *http.Request) {
	fmt.Fprint(writer, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:jmstudyacc@gmail.com\">jmstudyacc@gmail.com</a>")
}

// implementing a custom router
func pathHandler(writer http.ResponseWriter, reader *http.Request) {
	// implement a switch statement to capture the options
	switch reader.URL.Path {
	case "/":
		homeHandler(writer, reader)
	case "/contact":
		contactHandler(writer, reader)
	default:
		// this allows you to customise the response
		http.Error(writer, "Error 404\nPage Not Found - Did you mean something else?", http.StatusNotFound)
	}
}

func main() {
	fmt.Println("Starting the webserver on :3000")
	// starts the server & prevents the code from exiting
	err := http.ListenAndServe(":3000", http.HandlerFunc(pathHandler))

	if err != nil {
		panic(err)
	}
}
