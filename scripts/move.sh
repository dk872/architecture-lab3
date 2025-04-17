#!/bin/bash

startX=0.0
startY=0.0
step=0.05

x=$startX
y=$startY
dir=1
curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "update"
curl -X POST http://localhost:17000 -d "figure $startX $startY"
curl -X POST http://localhost:17000 -d "update"

while true; do
    reachedMax=$(echo "$x >= 1 || $y >= 1" | bc)
    reachedMin=$(echo "$x <= 0 || $y <= 0" | bc)

    if [[ $reachedMax -eq 1 ]]; then
        dir=-1
    elif [[ $reachedMin -eq 1 ]]; then
        dir=1
    fi

    x=$(echo "$x + $step*$dir" | bc)
    y=$(echo "$y + $step*$dir" | bc)

    curl -X POST http://localhost:17000 -d "move $x $y"
    curl -X POST http://localhost:17000 -d "update"

    sleep 1
done
