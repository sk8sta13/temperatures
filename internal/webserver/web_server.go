package webserver

import (
	"log"
	"net/http"

	"github.com.br/sk8sta13/temperatures/internal/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type HandlerProps struct {
	Method string
	Path   string
	Func   http.HandlerFunc
}

type WebServer struct {
	Router   chi.Router
	Handlers []HandlerProps
}

func NewWebServer() *WebServer {
	newWebServer := WebServer{
		Router:   chi.NewRouter(),
		Handlers: make([]HandlerProps, 0),
	}

	newWebServer.AddHandler(http.MethodGet, "/", handlers.ZipCodeAndTemperature)

	return &newWebServer
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, HandlerProps{
		Method: method,
		Path:   path,
		Func:   handler,
	})
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, h := range s.Handlers {
		s.Router.Method(h.Method, h.Path, h.Func)
	}

	if err := http.ListenAndServe("0.0.0.0:8080", s.Router); err != nil {
		log.Printf("Error starting the server.")
		return
	}
}
