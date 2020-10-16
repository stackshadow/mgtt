# MGTT

Another M(Go)TT-Broker written in Go

This MQTT-Broker is currently under development. But basic functionality works ( Publish/Subscribe with retained )

## Why

All the other broker, that i checked, are complicated to read/understand. 
This project aim's NOT to be feature-complete MQTT-Broker. It's should just be a simple broker with basic functionality.

## Features

- [x] Connect ( CONNECT / CONACK )
- [x] Ping ( PINGREQ/PINGRESP )
- [x] Publish ( PUBLISH / PUBACK )
- [x] Subscribe ( SUBSCRIBE / SUBACK )
- [ ] Local subsribe with callback
- [ ] Unsubscribe
- [x] QoS 0 messages
- [x] QoS 1 messages
- [ ] QoS 2 messages
- [x] Retained messages stored on [boltdb](https://github.com/boltdb/bolt) on disk
- [ ] Will messages
- [ ] Plugins ( Prepared, on going )
- [ ] TLS/SSL
- [ ] Disconnect

- [x] Zerolog with terminal-output and json-output-support
- [x] Kong command-line-parser with environment-support
- [ ] Dockerfile
- [ ] Docker-Compose

## Plugins
- [ ] Username/Password auth
- [ ] ACL
- [ ] Metrics
- [ ] $SYS-Support

