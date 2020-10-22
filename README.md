# MGTT

Another M(Go)TT-Broker written in Go

This MQTT-Broker is currently under development. But basic functionality works ( Publish/Subscribe with retained )

## Why

All the other broker, that i checked, are complicated to read/understand.
This broker should be easy and understandable with some extra sugar

And now i found out, that for QoS-2 the "specification" not suit my needs:
- The PUBCOMP is sended to the publisher when the broker publish the message to an subscriber
- For the Publisher, the protocol ends and the message is treated as received
- But i really want to know if the QoS-2-Message was delivered
- So MGTT is waiting that an client send back PUBREC

## Features

- [x] Connect ( CONNECT / CONACK )
- [x] Ping ( PINGREQ/PINGRESP )
- [x] Publish ( PUBLISH / PUBACK )
- [x] Subscribe ( SUBSCRIBE / SUBACK )
- [ ] Unsubscribe
- [x] QoS 0 messages
- [x] QoS 1 messages
- [x] QoS 2 messages
- [x] Retained messages stored on [boltdb](https://github.com/boltdb/bolt) on disk
- [ ] Will messages
- [ ] Plugins ( Prepared, on going )
- [x] TLS/SSL
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

