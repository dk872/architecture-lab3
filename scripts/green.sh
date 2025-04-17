#!/bin/bash

curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "update"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "bgrect 0.25 0.25 0.75 0.75"
curl -X POST http://localhost:17000 -d "update"