// tells Go - package main is the entrypoint and main() is where the code starts
package main

// importing fmt to print out stuff and net/http to set up webserver
import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// creating a function to handle the http request & present a web page
func homeHandler(writer http.ResponseWriter, reader *http.Request) {
	// defines the Content-Type header to be used >> Content Type Headers - HTTP Headers ~ Mozilla
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	// agnostic example - filepath.Join inserts the appropriate path
	tmplPath := filepath.Join("templates", "home.gohtml")
	// parse the template itself and then execute - UNIX based example = template.ParseFiles("templates/home.gohtml")
	templ, err := template.ParseFiles(tmplPath)
	if err != nil {
		// log = inclusion of timestamps etc. in the log messages
		log.Printf("Parsing template %v.", err)
		http.Error(writer, "There was an error parsing the template",
			http.StatusInternalServerError)
		// the return here tells the code to stop running
		return
	}

	err = templ.Execute(writer, nil)
	if err != nil {
		log.Printf("Executing template %v.", err)
		http.Error(writer, "There was an error executing the template",
			http.StatusInternalServerError)
		// the return here tells the code to stop running
		return
	}
}

func contactHandler(writer http.ResponseWriter, reader *http.Request) {
	fmt.Fprint(writer, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:jmstudyacc@gmail.com\">jmstudyacc@gmail.com</a></p>")

}

func myRequestHandler(writer http.ResponseWriter, reader *http.Request) {
	userID := chi.URLParam(reader, "userID")

	_, err := writer.Write([]byte(fmt.Sprintf("Hi there user %v", userID)))
	fmt.Printf("User has ID: %v\n", userID)
	if err != nil {
		return
	}

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

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/users/{userID}", myRequestHandler)
	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "Page Not Found", http.StatusNotFound)
	})
	fmt.Println("Starting the webserver on :3000")
	// starts the server & prevents the code from exiting
	err := http.ListenAndServe(":3000", r)

	if err != nil {
		panic(err)
	}
}
