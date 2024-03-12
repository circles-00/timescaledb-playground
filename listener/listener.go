package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"timeseries_db/mqtt"

	go_mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
)

func makeDBConnection() *sql.DB {
	// PostgreSQL connection parameters
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "password" // your password
	dbname := "db_name" // your database name

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Try to ping the database to ensure the connection is working
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func insertData(db *sql.DB, data []byte) {
	time := time.Now()

	stmt, err := db.Prepare("INSERT INTO data(time, data) VALUES ($1, $2)")

	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(time, data)

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted data into database")
}


func createMqttListener(db *sql.DB) {
	clientId := "go_mqtt_listener"
	client := mqtt.CreateMqttClient(clientId)

	topic := "test/topic"

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}

	// Define a message handler
	client.Subscribe(topic, 0, func(client go_mqtt.Client, msg go_mqtt.Message) {
		fmt.Printf("Received message on topic %s: %d\n", msg.Topic(), len(msg.Payload()))
    insertData(db, msg.Payload())
	})

	fmt.Printf("Subscribed to topic: %s\n", topic)

	// Setup a signal handler to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		client.Disconnect(250) // Disconnect with a timeout of 250 milliseconds
		os.Exit(0)
	}()

	// Wait indefinitely
	select {}
}

func main() {
	fmt.Println("Successfully connected to PostgreSQL!")

	db := makeDBConnection()

	defer db.Close()

	createMqttListener(db)
}
