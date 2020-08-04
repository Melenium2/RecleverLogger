package handlers

import (
	"encoding/json"
	"errors"
	"github.com/RecleverLogger/logger"
	"net/http"
)

type Handler struct {
	Path       string
	HandleFunc func(w http.ResponseWriter, r *http.Request)
	Method     string
}

type Handlers map[string]Handler

func (s *Service) Log(w http.ResponseWriter, r *http.Request) {
	log := &logger.SingleLog{}
	if err := json.NewDecoder(r.Body).Decode(log); err != nil {
		s.logger.Logf("Error while decoding json, %v", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if log.Type == "" {
		s.logger.Logf("Message type is empty")
		writeError(w, http.StatusBadRequest, errors.New("empty request"))
		return
	}

	if err := s.db.Save(r.Context(), log); err != nil {
		s.logger.Logf("Error from db while saving log, %v", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	s.logger.Logs(log)


	writeResponse(w, http.StatusOK, "Recorded")
}

func writeResponse(w http.ResponseWriter,code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			writeError(w, http.StatusInternalServerError, err)
		}
	}
}

func writeError(w http.ResponseWriter, code int, err error) {
	writeResponse(w, code, map[string]interface{} {
		"error": err.Error(),
	})
}