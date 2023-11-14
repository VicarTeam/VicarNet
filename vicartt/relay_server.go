package vicartt

import (
	"io"

	"github.com/gin-gonic/gin"
)

func NewServer() *Server {
	server := &Server{
		clients:       make(map[string]*Client),
		onClientOpen:  make(chan *Client),
		onClientClose: make(chan *Client),
	}

	go server.listen()

	return server
}

type Server struct {
	clients       map[string]*Client
	onClientOpen  chan *Client
	onClientClose chan *Client
}

func (s *Server) listen() {
	for {
		select {
		case client := <-s.onClientOpen:
			s.clients[client.AccessKey] = client
		case client := <-s.onClientClose:
			delete(s.clients, client.AccessKey)
			close(client.conn)
		}
	}
}

func (s *Server) HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

func (s *Server) PrepareMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &Client{
			AccessKey: c.Query("accessKey"),
			conn:      make(chan string),
		}

		s.onClientOpen <- client

		defer func() {
			s.onClientClose <- client
		}()

		c.Set("ttClient", client)
		c.Next()
	}
}

func (s *Server) HandleSSE() gin.HandlerFunc {
	return func(c *gin.Context) {
		client, ok := c.MustGet("ttClient").(*Client)
		if !ok {
			panic("invalid client")
		}

		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-client.conn; ok {
				c.SSEvent("message", msg)
				return true
			}

			return false
		})
	}
}

func (s *Server) RelayMessage(id string, msg interface{}) bool {
	client, ok := s.clients[id]
	if !ok {
		return false
	}

	client.Send(msg)

	return true
}
