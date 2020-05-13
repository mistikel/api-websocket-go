#!/bin/bash
for ((request=1;request<=20;request++))
do
     curl -H "Content-Type: application/json" -d '{"message":"publish messages"}' -X POST http://localhost:8080/message
done