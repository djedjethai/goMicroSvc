# goMicroSvc
Build a microservices architecture, each service being itself implemented following the hexagonal architecture.
A Broker facilitate the communication between services and an event bus(RabbitMQ here) reduces coupling between components by removing the need for identity on the
connector interface, which finally match an architecture of type Brokered Distributed Objects.
The objectif here was not to implement a fully functional application but to focus of the connections, the data transfert protocol and the messaging(amqp protocol) events bus. 

## Technologies
- Golang
- Postgres
- MongoDB
- Communicate between services using JSON
- Remote Procedure Calls
- gRPC
- Messaging protocol(amqp) using RabbitMQ
- Deployment using Docker-swarm
- Deployment using Kubernetes

## key issues
- Implementing each service following the hexagonal architecture
- Using RPC and/or gRPC to hasten the communication between services 
- Setting up RabbbitMQ to provide event-messaging communication between services  
