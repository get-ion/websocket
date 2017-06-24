package main

import (
	"fmt"

	"github.com/get-ion/ion"
	"github.com/get-ion/ion/context"

	"github.com/get-ion/websocket"
)

func main() {
	app := ion.New()

	app.Get("/", func(ctx context.Context) {
		ctx.ServeFile("websockets.html", false) // second parameter: enable gzip?
	})

	setupWebsocket(app)

	// x2
	// http://localhost:8080
	// http://localhost:8080
	// write something, press submit, see the result.
	app.Run(ion.Addr(":8080"))
}

func setupWebsocket(app *ion.Application) {
	// create our echo websocket server
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(handleConnection)

	// register the server on an endpoint.
	// see the inline javascript code i the websockets.html, this endpoint is used to connect to the server.
	app.Get("/echo", ws.Handler())

	// serve the javascript built'n client-side library,
	// see weboskcets.html script tags, this path is used.
	app.Any("/ion-ws.js", func(ctx context.Context) {
		ctx.Write(websocket.ClientSource)
	})
}

func handleConnection(c websocket.Connection) {
	// Read events from browser
	c.On("chat", func(msg string) {
		// Print the message to the console, c.Context() is the ion's http context.
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		// Write message back to the client message owner:
		// c.Emit("chat", msg)
		c.To(websocket.Broadcast).Emit("chat", msg)
	})
}
