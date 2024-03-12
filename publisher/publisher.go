package main

import (
	"fmt"
	"time"
	"timeseries_db/data"
	"timeseries_db/mqtt"
)

func createMqttPublisher(publisherId int) {
	for true {
		clientId := "go_mqtt_publisher" + fmt.Sprintf("%d", publisherId)
		client := mqtt.CreateMqttClient(clientId)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		topic := "test/topic"
		payload := data.Data

		// Publish a message
		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		fmt.Printf("%s: Published message on topic: %s\n", clientId, topic)

		time.Sleep(1 * time.Second)
	}
}

func main() {
	for i := 0; i < 100; i++ {
		go func(publisherId int) {
			fmt.Printf("Created a publisher #%d\n", publisherId)
			createMqttPublisher(publisherId)
		}(i)

		time.Sleep(2 * time.Second)
	}
}
