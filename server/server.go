package server

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

// Server is a simple web server.
type Server struct {
	bindAddr string
	config   *Config
	secret   []byte
	logger   *log.Logger
}

// NewServer creates a `Server` instance and register request handlers
func NewServer(bindAddr string, config *Config) (*Server, error) {
	s := &Server{
		bindAddr: bindAddr,
		config:   config,
	}

	logFile, err := os.OpenFile(s.config.LogPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}

	s.logger = log.New(logFile, "", log.LstdFlags)

	secret, err := loadSecret(s.config.SecretPath)
	if err != nil {
		s.logger.Printf("Warning: cannot load secret key at path %q: %v", s.config.SecretPath, err)
		s.logger.Print("Please make sure the secret key file contains 16, 24, or 32 bytes " +
			"of key data encoded in base64 format.")
		s.logger.Printf("You can generate a base64 key data from http://%s/genkey.", bindAddr)
	} else {
		s.logger.Printf("Loaded secret key from %q", s.config.SecretPath)
		s.secret = secret
	}

	http.HandleFunc("/headers", s.headersHandler)
	http.HandleFunc("/env", s.envHandler)
	http.HandleFunc("/encrypt", s.encryptHandler)
	http.HandleFunc("/decrypt/", s.decryptHandler)
	http.HandleFunc("/genkey", s.genKeyHandler)
	http.HandleFunc("/", s.indexHandler)
	return s, nil
}

// ListenAndServe listens and serves HTTP requests.
func (s *Server) ListenAndServe() error {
	log.Printf("Listening on %s", s.bindAddr)
	return http.ListenAndServe(s.bindAddr, nil)
}

func (s *Server) indexHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "echo")
}

func (s *Server) headersHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r.Header)
}

func (s *Server) envHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, os.Environ())
}

func (s *Server) ensureSecret() error {
	if s.secret != nil {
		return nil
	}
	log.Printf("Bad secret key configuration. See %q for more detail.", s.config.LogPath)
	return errors.New("Bad server configuration")
}

func (s *Server) encryptHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.ensureSecret(); err != nil {
		writeError(w, err)
		return
	}
	encrypted, err := encrypt("super secret message", s.secret)
	if err != nil {
		writeError(w, err)
		return
	}

	writeBinary(w, encrypted)
}

func (s *Server) decryptHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.ensureSecret(); err != nil {
		writeError(w, err)
		return
	}
	ciphertext := path.Base(r.URL.Path)
	decrypted, err := decrypt(ciphertext, s.secret)
	if err != nil {
		writeError(w, err)
		return
	}

	fmt.Fprint(w, decrypted)
}

func (s *Server) genKeyHandler(w http.ResponseWriter, _ *http.Request) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		writeError(w, err)
		return
	}

	writeBinary(w, key)
}
