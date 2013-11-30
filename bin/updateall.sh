#!/bin/sh

set -x # echo on

echo "Updating all dependencies..."

go get -u github.com/axgle/pinyin
go get -u github.com/go-sql-driver/mysql

go get -u github.com/gorilla/context
go get -u github.com/gorilla/mux
go get -u github.com/gorilla/schema

go get -u github.com/elivoa/gxl
go get -u github.com/elivoa/got
go get -u github.com/elivoa/syd

echo "All Done!"
