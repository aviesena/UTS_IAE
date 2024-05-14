package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/mux"
    "github.com/streadway/amqp"
)

type MessageContent struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Email       string `json:"email"`
    PhoneNumber string `json:"phone_number"`
    Address     string `json:"address"`
}

type Message struct {
    Message  MessageContent `json:"message"`
    DateTime string         `json:"dateTime"`
}

var messagesQueue []Message
var mutex sync.Mutex

func main() {
    // Connect to RabbitMQ server
    conn, err := amqp.Dial("amqp://localhost")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    // Create a new channel
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    // Create the exchange
    err = ch.ExchangeDeclare(
        "userExchange", // name
        "fanout",       // type
        true,           // durable
        false,          // auto-deleted
        false,          // internal
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare an exchange: %v", err)
    }

    // Create the queue
    q, err := ch.QueueDeclare(
        "orderQueue", // name
        true,         // durable
        false,        // delete when unused
        false,        // exclusive
        false,        // no-wait
        nil,          // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }

    // Bind the queue to the exchange
    err = ch.QueueBind(
        q.Name,       // queue name
        "",           // routing key
        "userExchange", // exchange
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to bind a queue: %v", err)
    }

    // Consume messages from the queue
    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    go func() {
        for d := range msgs {
            var msg Message
            err := json.Unmarshal(d.Body, &msg)
            if err != nil {
                continue
            }

            mutex.Lock()
            messagesQueue = append(messagesQueue, msg)
            mutex.Unlock()
        }
    }()

    // Set up HTTP server
    r := mux.NewRouter()
    r.HandleFunc("/", getMessages).Methods("GET")

    log.Printf("Server started on port %v", 3001)
    log.Fatal(http.ListenAndServe(":3001", r))
}

func getMessages(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    defer mutex.Unlock()

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(messagesQueue); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
