package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

func CreateMqttClient(clientId string) mqtt.Client {
	// MQTT broker URL. Change this to your MQTT broker's URL.
	broker := "tcp://localhost:1883"

	// Create an MQTT client options
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientId)

	client := mqtt.NewClient(opts)

	return client
}


