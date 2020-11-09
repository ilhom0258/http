package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

//HandleFunc handling functions
type HandleFunc func(conn net.Conn)

//Server structure
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandleFunc
}

//const hosts for tests
const (
	Host = "0.0.0.0"
	Port = "9999"
)

//NewServer Creates new server
func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		handlers: make(map[string]HandleFunc),
	}
}

//Register registers handlers 
func (s *Server) Register(path string, handler HandleFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	//TODO start server on host & port
	listner, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
	}

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
		log.Print("requestLineEndErr: ", requestLineEnd)
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		log.Print("partsErr: ", parts)
	}

	s.mu.RLock()
	if handler, ok := s.handlers[parts[1]]; ok {
		s.mu.RUnlock()
		handler(conn)
	}
	return
}

func (s *Server) generateResponse(body string) string {
	return "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
		"Content-Type: text/html\r\n" +
		"Connection: close\r\n" +
		"\r\n" + body
}

func (s *Server) RouteHandler(body string) func(conn net.Conn) {
	return func(conn net.Conn) {
		_, err := conn.Write([]byte(s.generateResponse(body)))
		if err != nil {
			log.Print(err)
		}
	}
}
