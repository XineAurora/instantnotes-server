package instantnotesapi

import "github.com/XineAurora/instantnotes-server/internal/handler"

type APIServer struct {
	Handler *handler.Handler
}

func New() *APIServer {
	apiserver := APIServer{Handler: handler.New()}
	return &apiserver
}

func (s *APIServer) Start() error {
	if err := s.Handler.Router.Run(); err != nil {
		return err
	}
	return nil
}
