# MGTT

Another M(Go)TT-Broker written in Go

THIS PROJECT IS MAINTAINED ON GITLAB.COM ( https://gitlab.com/stackshadow/mgtt ) Please create feature requests or bugs there :)

## Usage

Read the docs in the ./docs folder

## Capabilitys

- [x] Connect ( CONNECT / CONACK )
- [x] Ping ( PINGREQ/PINGRESP )
- [x] Publish ( PUBLISH / PUBACK )
- [x] Subscribe ( SUBSCRIBE / SUBACK )
- [x] Unsubscribe
- [x] QoS 0 messages
- [x] QoS 1 messages
- [x] QoS 2 messages
- [x] Retained messages stored on [boltdb](https://github.com/boltdb/bolt) on disk
- [x] Automatic resending of failed packets
- [x] Will message
- [x] Plugins
- [x] TLS/SSL
- [x] Disconnect
- [x] Sessions

## Features

- [x] Zerolog with terminal-output and json-output-support
- [x] Kong command-line-parser with environment-support
- [x] Dockerfile
- [ ] Healthcheck
- [x] Docker-Compose
- [ ] Build with buildah
- [x] $SYS-Support

## Plugins

- [x] Username/Password auth
- [x] ACL
- [ ] Metrics
