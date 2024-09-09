package webserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      map[string]http.HandlerFunc{},
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	for methodPath, handler := range s.Handlers {
		parts := strings.SplitN(methodPath, ":", 2)
		if len(parts) != 2 {
			panic("Invalid method format required: 'METHOD:path' " + methodPath)
		}
		method := strings.ToUpper(parts[0])
		path := parts[1]

		switch method {
		case "GET":
			s.Router.Get(path, handler)
		case "POST":
			s.Router.Post(path, handler)
		case "PUT":
			s.Router.Put(path, handler)
		case "DELETE":
			s.Router.Delete(path, handler)
		default:
			panic("Invalid method: " + method)
		}
	}
	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		panic(err)
	}
	return nil
}
