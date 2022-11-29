package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	guuid "github.com/google/uuid"
	"net/http"
	"video-chat/pkg/webrtc"
)

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func RoomCreate(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("/room/%s", guuid.NewString()))
}

func Room(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	if uuid == "" {
		c.Status(http.StatusBadRequest)

		return nil
	}

	uuid, suuid, _ := createOrGetRoom(uuid)
}

func RoomWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")

	if uuid == "" {
		return
	}

	_, _, room := createOrGetRoom(uuid)

	webrtc.RoomConn(room.Peers)
}

func createOrGetRoom(uuid string) (string, string, webrtc.Room) {

}

func RoomViewerWebSocket(c *websocket.Conn) {

}

func roomViewerConn(c *websocket.Conn, p *webrtc.Peers) {

}
