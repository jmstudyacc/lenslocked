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
	fmt.Fprint(writer, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:jmstudyacc@gmail.com\">jmstudyacc@gmail.com</a></p>")
}

// creating a function to handle the HTTP request & response for FAQ page
func faqHandler(writer http.ResponseWriter, reader *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf8")
	fmt.Fprint(writer, `<h1>FAQ Page</h1>
<ul>
	<li><h3>Is there a Free Version?</h3><p>Yes! We offer a free trial for 30 days on any paid plans.</p>
	<li><h3>What are your Support hours?</h3><p>We have support staff answering emails 24/7, though response times may be slower at weekends,</p>
	<li><h3>How do I contact Support?</h3><p>Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a></p>
</ul>`)
}

// implementing a custom router
//func pathHandler(writer http.ResponseWriter, reader *http.Request) {
// implement a switch statement to capture the options
//	switch reader.URL.Path {
//	case "/":
//		homeHandler(writer, reader)
//	case "/contact":
//		contactHandler(writer, reader)
//	default:
// this allows you to customise the response
//		http.Error(writer, "Error 404\nPage Not Found - Did you mean something else?", http.StatusNotFound)
//	}
//}

type Router struct{}

// implementing a custom router
func (Router) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	// implement a switch statement to capture the options
	switch reader.URL.Path {
	case "/":
		homeHandler(writer, reader)
	case "/contact":
		contactHandler(writer, reader)
	case "/faq":
		faqHandler(writer, reader)
	default:
		// this allows you to customise the response
		writer.Header().Set("Content-Type", "text/html;charset=utf8")
		http.Error(writer, "Error 404\nPage Not Found - Did you mean something else?", http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the webserver on :3000")
	// starts the server & prevents the code from exiting
	err := http.ListenAndServe(":3000", router)

	if err != nil {
		panic(err)
	}
}
