package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mistikel/api-websocket-go/pubsub"
	"github.com/mistikel/api-websocket-go/utils"
	uuid "github.com/satori/go.uuid"
)

var IsShuttingDown = false

type Handler interface {
	CreateRouter() *mux.Router
	HealthCheck(w http.ResponseWriter, r *http.Request)
	HomePage(w http.ResponseWriter, r *http.Request)
	StoreMessage(w http.ResponseWriter, r *http.Request)
	ResolveMessage(w http.ResponseWriter, r *http.Request)
	SubscribeMessage(w http.ResponseWriter, r *http.Request)
}

func NewHandler() Handler {

	return &Controller{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		pubsub:  pubsub.New(),
		Message: make(chan string, 1),
	}
}

type M struct {
	Message string `json:"message"`
}

type Controller struct {
	upgrader websocket.Upgrader
	pubsub   pubsub.Manager

	Message chan string
}

func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if IsShuttingDown {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Preparing to shutdown"))
		return
	}
	response := "OK"
	status := http.StatusOK
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func (c *Controller) CreateRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", c.HomePage).Methods("GET")
	r.HandleFunc("/message", c.ResolveMessage).Methods("GET")
	r.HandleFunc("/message", c.StoreMessage).Methods("POST")
	r.HandleFunc("/subscribe", c.SubscribeMessage)
	return r
}

func (c *Controller) HomePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

func (c *Controller) StoreMessage(w http.ResponseWriter, r *http.Request) {
	var m M
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		response := fmt.Sprintf("Error: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(response))
		return
	}

	c.pubsub.Publish(r.Context(), m.Message)

	response := "Message well received, have a nice day!"
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(response))
}

func (c *Controller) ResolveMessage(w http.ResponseWriter, r *http.Request) {
	messages := c.pubsub.ResolveAllMessages(r.Context())
	js, err := json.Marshal(messages)
	if err != nil {
		response := fmt.Sprintf("Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response))
		return
	}
	w.Write(js)
}

func (c *Controller) SubscribeMessage(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		response := fmt.Sprintf("Error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response))
		return
	}
	id := uuid.NewV4()

	pingTicker := time.NewTicker(time.Minute)
	defer func() {
		conn.Close()
		pingTicker.Stop()
		c.pubsub.Unsubscribe(r.Context(), id.String())
	}()

	closedChan := make(chan bool, 0)
	go func() {
		for {
			mType, _, err := conn.ReadMessage()
			if mType == websocket.CloseMessage || err != nil {
				closedChan <- true
				return
			}
		}
	}()

	sChan := c.pubsub.Subscribe(r.Context(), id.String())
	for {
		select {
		case msg := <-sChan:
			utils.InfoContext(r.Context(), "Message: %v", msg)
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg+id.String()))
			if err != nil {
				utils.ErrContext(r.Context(), "[Subcriber]: %v", err)
				return
			}
		case <-pingTicker.C:
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				utils.ErrContext(r.Context(), "[Subcriber]: %v", err)
				return
			}
		case <-closedChan:
			utils.InfoContext(r.Context(), "Closing websocket connection")
			return
		}
	}
}
