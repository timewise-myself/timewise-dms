package main

import (
	_ "dbms/docs"
	"dbms/server"
)

// @title timewise-dbms
// @version 1.0
// @description Timewise database management system
// @in header
func main() {
	server.RegisterServer()
}
