#!/bin/bash

curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "bgrect 0.18 0.18 0.82 0.82"
curl -X POST http://localhost:17000 -d "figure 0.5 0.5"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "figure 0.55 0.55"
curl -X POST http://localhost:17000 -d "update"
