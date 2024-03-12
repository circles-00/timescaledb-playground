# Playground for [TimescaleDB](https://github.com/timescale/timescaledb)

Using mosquitto broker for communication.
The Golang part includes a listener that listens on a topic and publishers that will run in parallel and send out data to the same topic that the listener is subscribed to.

The listener inserts the data into the database.
