package http

import (
	"errors"
	"fmt"
	"furnace-monitoring-system-server/pkg/app"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	ADDRESS = ":5000"
)

type Server struct {
	endpointHandler *EndpointHandler
	assetsHandler   *EmbedadAssetsHandler
	router          *chi.Mux
}

func CreateServer(measurementService *app.MeasurementService) (*Server, error) {
	if measurementService == nil {
		return nil, errors.New("Server: Provided MeasurementService reference in nil")
	}

	server := new(Server)
	server.assetsHandler = CreateEmbededAssetsHandler()

	endpointHandler, err := CreateHandler(measurementService, server.assetsHandler)
	if err != nil {
		return nil, fmt.Errorf("Server: Failed to create Handler instance: %w", err)
	}

	server.endpointHandler = endpointHandler
	server.router = chi.NewRouter()

	// TODO: Logger middleware
	server.router.Use(middleware.NoCache)
	server.router.Use(middleware.RealIP)
	server.router.Use(middleware.AllowContentEncoding("gzip", "deflate"))
	// TODO: Append Content-Type to all handler functions
	server.router.Use(middleware.Compress(5, "text/html", "text/css"))
	server.router.Use(middleware.CleanPath)
	server.router.Use(middleware.Recoverer)

	server.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	assetsHandler := server.assetsHandler.GetEmbededStaticAssetsHandler()
	server.router.Handle("/assets/", http.StripPrefix("/assets", assetsHandler))

	server.router.Get("/", server.endpointHandler.HandleIndexTemplate)
	server.router.Get("/error", server.endpointHandler.HandleErrorTemplate)

	server.router.HandleFunc("/sensor", server.endpointHandler.HandleSensorSocket)

	return server, nil
}

func (srv *Server) Listen() error {
	return http.ListenAndServe(ADDRESS, srv.router)
}
