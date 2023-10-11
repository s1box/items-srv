/*
(c) 2023 Mykola Morhun
Demo GoLang REST API server with MySQL database.

This program is free software: you can redistribute it and/or modify it under the terms of the
GNU General Public License as published by the Free Software Foundation,
either version 2 of the License, or (at your option) any later version.
This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
See the GNU General Public License for more details.
*/
package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/s1box/dc-go-srv/api"
)

func main() {
	srvHost := "localhost"
	srvPort := "8080"
	if host := os.Getenv("HOSTNAME"); host != "" {
		srvHost = host
	}
	if port := os.Getenv("PORT"); port != "" {
		srvPort = port
	}

	runServer(srvHost, srvPort)
}

func runServer(host, port string) {
	itemsApi := api.NewItemsRestApi()
	statusApi := api.NewStatusRestApi()

	router := gin.Default()

	itemsApi.AddApiToRouter(router)
	statusApi.AddApiToRouter(router)

	router.Run(fmt.Sprintf("%s:%s", host, port))
}
