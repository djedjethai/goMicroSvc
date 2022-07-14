#!/bin/bash

curl --header "Content-Type:application/json" --request POST http://localhost:8080 \ 

curl --header "Content-Type:application/json" --request POST --data '{"action":"log", "log":{"name":"event","data":"grpc datas"}}' http://localhost:8080/log-grpc \

curl --header "Content-Type:application/json" --request POST --data '{"action":"log", "log":{"name":"event","data":"grpc datas"}}' http://localhost:8080/handle \


curl --header "Content-Type:application/json" --request POST --data '{"action":"mail", "mail":{"from":"me@exemple.com","to":"you@there.com","subject":"test email","message":"hello man"}}' http://localhost:8080/handle \


curl --header "Content-Type:application/json" --request POST --data '{"action":"auth","auth":{"email":"admin@example.com","password":"verysecret"}}' http://localhost:8080/handle \


