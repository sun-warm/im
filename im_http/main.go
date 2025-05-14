package main

import (
	"im_http/im_ws_client"
	"im_http/server"
)

func main() {
	im_ws_client.StartIMWSClient()
	server.StartIMHTTPServer()
}
