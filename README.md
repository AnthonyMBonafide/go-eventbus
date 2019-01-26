[![Build Status](https://travis-ci.com/AnthonyMBonafide/go-eventbus.svg?branch=master)](https://travis-ci.com/AnthonyMBonafide/go-eventbus)
[![Go Report Card](https://goreportcard.com/badge/github.com/AnthonyMBonafide/go-eventbus)](https://goreportcard.com/report/github.com/AnthonyMBonafide/go-eventbus)
[![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)](https://github.com/AnthonyMBonafide/go-eventbus/blob/initial-work/LICENSE)

# go-eventbus
Event bus for the GO programming language used to facilitate communication
between clustered GO applications either in distributed or non-distributed
systems. This tool is inspired by the Vert.X eventbus.

# Features

The go-eventbus provides the following functionality:

- Distributed caching
- Point-to-point message sending via topic
- Broadcasting messages via topic
- Registering consumers to a topic
- Message filtering by properties
- Auto discovery of cluster nodes
- Load balancing
- Metrics
- Logging
- Timer tasks

# Initial Tasks

1. Determine a communication protocol used between nodes. Done, Websockets
1. Get multiple nodes to talk to each other, no discovery
1. Implement a heartbeat / remove bad nodes from cluster
1. Add quorum to prevent split brain
1. Add support for registering consumers of a topic
1. Add support for sending messages to a topic (Publish and/or point-to-point)
1. Add support for message delivery
1. Add support for message reply to a temp topic
