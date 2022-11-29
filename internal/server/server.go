package server

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
	"os"
	"time"
	"video-chat/internal/handlers"
	"video-chat/pkg/webrtc"
)

var (
	addr = flag.String("addr", ":"+os.Getenv("PORT"), "")
	cert = flag.String("cert", "", "")
	key  = flag.String("key", "", "")
)

func Run() error {
	flag.Parse()

	if *addr == ":" {
		*addr = ":8080"
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", handlers.Welcome)

	app.Get("/room/create", handlers.RoomCreate)
	app.Get("/room/:uuid", handlers.Room)
	app.Get("/room/:uuid/websocket", websocket.New(handlers.RoomWebsocket, websocket.Config{HandshakeTimeout: 10 * time.Second}))
	app.Get("/room/:uuid/chat", handlers.RoomChat)
	app.Get("/room/:uuid/chat/websocket", websocket.New(handlers.RoomChatWebsocket))
	app.Get("/room/:uuid/viewer/websocket", websocket.New(handlers.RoomViewerWebsocket))

	app.Get("/stream/:ssuuid", handlers.Stream)
	app.Get("/stream/:ssuuid/websocket", websocket.New(handlers.StreamWebsocket, websocket.Config{HandshakeTimeout: 10 * time.Second}))
	app.Get("/stream/:ssuuid/chat/websocket", websocket.New(handlers.StreamChatWebsocket))
	app.Get("/stream/:ssuuid/viewer/websocket", websocket.New(handlers.StreamViewerWebsocket))

	app.Static("/", "./assets")

	webrtc.Rooms = make(map[string]*webrtc.Room)
	webrtc.Streams = make(map[string]*webrtc.Room)

	go dispatchKeyFrames()

	if *cert != "" {
		return app.ListenTLS(*addr, *cert, *key)
	}

	return app.Listen(*addr)
}

func dispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range webrtc.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}
