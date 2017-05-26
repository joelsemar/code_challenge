package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Service struct {
	dbConnection *MongoConnection
	router       *mux.Router
}

func (this *Service) serve(port int) {
	Log("Listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), this.router)
}

func (this *Service) setRoutes() {
	this.router = mux.NewRouter()
	this.router.HandleFunc("/chat", this.createMessage).Methods("POST")
	this.router.HandleFunc("/chat/{username}", this.getMessages).Methods("GET")
}

// create a new message according to the provided payload
func (this *Service) createMessage(responseWriter http.ResponseWriter, request *http.Request) {
	Log("Creating a message.. ")
	message := new(MessageDocument)

	if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
		this.send(responseWriter, ResponseMsg("Bad Request"), 400)
		log.Fatal(err)
		return
	}
	this.dbConnection.addMessage(message)
	this.created(responseWriter, message)

}

// Fetch all messages for the given username
func (this *Service) getMessages(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Get messages.. ")
	username := mux.Vars(request)["username"]
	results, err := this.dbConnection.getMessages(username)
	if err != nil {
		this.fail(responseWriter, err)
	}
	this.ok(responseWriter, results)
}

func (this *Service) ok(responseWriter http.ResponseWriter, response interface{}) {
	this.send(responseWriter, response, 200)
}

func (this *Service) created(responseWriter http.ResponseWriter, response interface{}) {
	this.send(responseWriter, response, 201)
}

func (this *Service) fail(responseWriter http.ResponseWriter, err error) {
	this.send(responseWriter, ResponseMsg(err.Error()), 500)
}

// Write a message out to the response writer, set headers and status code, marshall out given response object
func (this *Service) send(responseWriter http.ResponseWriter, response interface{}, statusCode int) {
	out, _ := json.Marshal(response)
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	fmt.Fprintf(responseWriter, "%s", out)
}

func NewService() (service *Service) {
	service = new(Service)
	service.dbConnection = NewMongoConnection()
	service.setRoutes()
	return
}
