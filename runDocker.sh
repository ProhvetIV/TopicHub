#!/bin/bash

docker build -t frontend ./vue-frontend
docker build -t backend ./
docker run -d -p 5173:80 frontend
docker run -d -p 8080:8080 backend
