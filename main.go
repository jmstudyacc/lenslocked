// tells Go - package main is the entrypoint and main() is where the code starts
//
// This branch shows the use of the go-chi/chi logger middleware
// Docs on this can be found at: https://github.com/go-chi/chi/tree/master/_examples/logging
package main

// importing fmt to print out stuff and net/http to set up webserver
import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
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

// Structured Logger is a simple, but powerful implementation of a custom structured
// logger back on logrus. Check out https://github.com/go-chi/httplog for dedicated pkg
// based on this work, designed for context-based http routers

type StructuredLogger struct {
	Logger *logrus.Logger
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("Request Complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func (l *StructuredLogger) NewLogEntry(reader *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(reader.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if reader.TLS != nil {
		scheme = "https"
	}

	logFields["http_scheme"] = scheme
	logFields["http_proto"] = reader.Proto
	logFields["http_method"] = reader.Method

	logFields["remote_addr"] = reader.RemoteAddr
	logFields["user_agent"] = reader.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, reader.Host, reader.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("Request Started")

	return entry
}

// Following are helper methods used by the application to get the
// request-scoped logger entry & set additional fields between handlers
//
// A useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .Info(), etc.

func GetLogEntry(reader *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(reader).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(reader *http.Request, key string, value interface{}) {
	if entry, ok := reader.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(reader *http.Request, fields map[string]interface{}) {
	if entry, ok := reader.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}

func main() {
	// sets up the logger backend using sirupsen/logrus and configure it
	// to use a custom JSONFormatter. Refer to the logrus docs on how to
	// configure the backend at https://github.com/sirupsen/logrus
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(NewStructuredLogger(logger))
	r.Use(middleware.Recoverer)

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
