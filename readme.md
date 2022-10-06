# Tracker decoder

Simple tracker decoder service, it recieves connections from devices and publishes the decoded packets into a rabbitMQ exchange

![image info](./docs/service_simplified.jpg)

---

## Configuration

configuration is set by a yml config file and enviroment variables, each variable on the yaml file can be overwrittern by a env var,
check `config/config.yml` for details.

when developing its easier to use the yml equivalent of those variables on your `config/config.dev.yml` file and running the service
with `make run_dev` or `go run cmd/main.go --config-file="./config/config.dev.yml"`


---

## Tracker configuration

The first thing you need is to get a public IP, chances are you will need to learn about [port forwarding](https://en.wikipedia.org/wiki/Port_forwarding).
if you want to easily debug connections that install [ngrok](https://ngrok.com/) and run `ngrok tcp <port_you_configured_on_port_forwarding>`

To confirm your ip address is open to the world send a tcp packet to your ip (or the ngrok ip if youre using it) and check for any logs.
Now you only need to configure your tracker to connect to your IP, check your tracker manual.

---

## Enviroment variables

|           name          |                                    meaning                                   | example                           |
|-------------------------|------------------------------------------------------------------------------|-----------------------------------|
| APP_DEBUG               | debug mode, if true will log to debug info to stdout                         | false                             |
| MAX_INVALID_PACKETS     | amount of invalid packets a tracker can send before the connection is closed | 20                                |
| RMQ_URL                 | rabbitmq url                                                                 | amqp://guest:guest@localhost:5672 |
| RMQ_EXCHANGE            | name of the rabbitmq exchange to publish events on                           | tracker_events_topic              |
| RMQ_RECONNECT_WAIT_TIME | seconds to wait before trying to reconnect when rabbitmq connection is lost  | 5                                 |
| TRACER_URL              | jaeger endpoint to send traces to                                            | http://localhost:14268/api/traces |
| TRACER_SERVICE_NAME     | name of the service to jaeger                                                | tracker_reciever                  |

---

## Rabbitmq

Tracker events are published to a single topic exchange with the following pattern <protocol_slug>.<event_type>.<device_imei>. eg: h02.location.imei. With this topic exchange your services can listen to events only they care about, examples:

- all event types of the h02 protocol `h02.*.*`
- location events regardless of the protocol and imei `*.location.*`
- events of a specific tracker, by its imei `*.*.8603412412412`

---

## Supported protocols / messages

- [h02](protocol/h02/doc/events.md)

---