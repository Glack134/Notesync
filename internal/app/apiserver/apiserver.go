package apiserver

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/polyk005/Notesync1.0/internal/app/blocknot"
	"github.com/sirupsen/logrus"
)

// APIServer
type APIServer struct {
	config   *Config
	logger   *logrus.Logger
	router   *mux.Router
	blocknot *blocknot.Blocknot
}

// New
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureBlocknot(); err != nil {
		return fmt.Errorf("failed to configure store: %w", err)
	}

	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.BinAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/main", s.handleHello())
}

func (s *APIServer) configureBlocknot() error {
	bl := blocknot.New(nil)
	if err := bl.Open(); err != nil {
		return err
	}
	s.blocknot = bl

	return nil
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
