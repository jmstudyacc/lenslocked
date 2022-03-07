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

func executeTemplate(wr http.ResponseWriter, filepath string) {
	wr.Header().Set("Content-Type", "text/html; charset=utf-8")
	// parse the template itself and then execute - UNIX based example = template.ParseFiles("templates/home.gohtml")
	templ, err := template.ParseFiles(filepath)
	if err != nil {
		// log = inclusion of timestamps etc. in the log messages
		log.Printf("Parsing template %v.", err)
		http.Error(wr, "There was an error parsing the template",
			http.StatusInternalServerError)
		// the return here tells the code to stop running
		return
	}

	err = templ.Execute(wr, nil)
	if err != nil {
		log.Printf("Executing template %v.", err)
		http.Error(wr, "There was an error executing the template",
			http.StatusInternalServerError)
		// the return here tells the code to stop running
		return
	}
}

// creating a function to handle the http request & present a web page
func homeHandler(writer http.ResponseWriter, reader *http.Request) {
	tmplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(writer, tmplPath)
}

func contactHandler(writer http.ResponseWriter, reader *http.Request) {
	tmplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(writer, tmplPath)
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
	// not creating a separate variable to hold the String - just placing into the function directly
	executeTemplate(writer, filepath.Join("templates", "faq.gohtml"))
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
