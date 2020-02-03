#!/bin/bash

curl -X POST http://0.0.0.0:8080/images/upload/$1 \
  -F "file=@$2" \
  -H "Content-Type: multipart/form-data"
