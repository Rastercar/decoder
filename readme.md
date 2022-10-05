# Tracker decoder

Simple tracker decoder service, it recieves connections from devices and publishes the decoded packets into a rabbitMQ exchange

![image info](./docs/service_simplified.jpg)

## Configuration

configuration is set by a yml config file and enviroment variables, each variable on the yaml file can be overwrittern by a env var,
check `config/config.yml` for details.

when developing its easier to use the yml equivalent of those variables on your `config/config.dev.yml` file and running the service
with `make run_dev` or `go run cmd/main.go --config-file="./config/config.dev.yml"`

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

## Supported protocols / messages

- h02
    - *Uplink*
    - ✅ heartbeat
    - ✅ real time location
    - ❌ location request
    - ❌ blind spots uploading
    - ❌ device alarm 
    - *Downlink*
    - ❌ cut-off/recover oil and engine 
    - ❌ response to location request
    - ❌ fortification (SF)
    - ❌ fortification (SF2)
    - ❌ disarming (CF)
    - ❌ disarming (CF2)
    - ❌ main number bind (UR)
    - ❌ server setting (IP)
    - ❌ terminal password setting
    - ❌ interval settings
    - ❌ allarm setting
    - ❌ device reboot
    - ❌ reset to defaults
    - ❌ network access point
    - ❌ answer mode
    - ❌ IMEI setting
    - ❌ language setting
    - ❌ audiomonitor
    - ❌ query device information
    - ❌ working mode setting