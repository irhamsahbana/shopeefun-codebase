package websocket

import (
	"codebase-app/internal/adapter"
	"codebase-app/pkg/errmsg"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func GroupChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomId := r.URL.Query().Get("room_id")
	if roomId == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp := `
			{
				"success": false,
				"message": "Room id is required",
				"data": null
			}
		`
		_, err := w.Write([]byte(resp))
		if err != nil {
			log.Error().Err(err).Msg("Error while writing response")
			return
		}
		return
	}

	claims, ok := r.Context().Value("claims").(map[string]any)
	if !ok {
		log.Error().Msg("Error while getting claims")
		return
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		log.Error().Msg("Error while getting user_id")
		return
	}

	role, ok := claims["role"].(string)
	if !ok {
		log.Error().Msg("Error while getting role")
		return
	}

	type User struct {
		UserId string `json:"user_id"`
		Name   string `json:"name"`
		Role   string `json:"role"`
	}

	var (
		db   = adapter.Adapters.ShopeefunPostgres
		user User
	)
	user.UserId = userId
	user.Role = role

	query := `SELECT name FROM users WHERE id = ?`
	err := db.QueryRowxContext(r.Context(), db.Rebind(query), userId).Scan(&user.Name)
	if err != nil {
		log.Error().Err(err).Any("user_id", userId).Msg("Error while getting user")
		return
	}

	ServeWs(rooms, roomId, user.UserId, user.Name, user.Role, w, r)
}

func ServeWs(rooms map[string]*Hub, roomId, userId, role, name string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading connection")
		return
	}

	room, exists := rooms[roomId]
	if !exists {
		room = NewHub(roomId)
		rooms[roomId] = room
		go room.Run(rooms)
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		room: room,
		id:   userId,
		role: role,
		name: name,
	}
	room.register <- client

	client.room.broadcast <- []byte(`
		{
			"event": "chat-group-user-join",
			"data": {
				"user_id": "` + userId + `",
				"name": "` + name + `",
				"message": "joined the chat"
			}
		}
	`)

	// Start reading and writing to the client
	go client.writePump()
	go client.readPump()
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		log.Info().Msg("Closing connection")
		c.room.unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error while closing connection")
		}
	}()

	var (
		v = adapter.Adapters.Validator
	)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Any("user_id", c.id).Msg("Error while reading message")
			break
		}

		// unmarshal the message to Event struct, so we can validate the event type
		var event Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Error().Err(err).Any("user_id", c.id).Msg("Error while unmarshalling message")
			c.send <- []byte(ERROR_PARSE_MSG)
			continue
		}

		// validate the event type, if it's invalid, send an error message to the client
		err = v.Validate(&event)
		_, errs := errmsg.Errors(err, &event)

		if err != nil {
			log.Error().Err(err).Any("user_id", c.id).Msg("Error while validating message")
			// marshal the invalid event to json
			invalidStructure := EventInvalid{
				Event:   VALIDATION_ERROR,
				Message: VALIDATION_ERROR_MESSAGE,
				Errors:  errs,
			}

			invalidStructureBytes, err := json.Marshal(&invalidStructure)
			if err != nil {
				log.Error().Err(err).Any("user_id", c.id).Msg("Error while marshalling invalid event")
				c.send <- []byte(ERROR_PARSE_MSG)
				continue
			}

			c.send <- invalidStructureBytes
			continue
		}

		switch event.Event {
		case CHAT_GROUP_USER_SEND_MSG:
			var payload EventUserSendMessage
			err = json.Unmarshal(message, &payload)
			if err != nil {
				log.Error().Err(err).Any("user_id", c.id).Msg("Error while unmarshalling message")
				c.send <- []byte(ERROR_PARSE_MSG)
				continue
			}

			err = v.Validate(&payload)
			_, errs = errmsg.Errors(err, &payload)
			if err != nil {
				log.Error().Err(err).Any("user_id", c.id).Msg("Error while validating message")
				invalidStructure := EventInvalid{
					Event:   VALIDATION_ERROR,
					Message: VALIDATION_ERROR_MESSAGE,
					Errors:  errs,
				}

				invalidStructureBytes, err := json.Marshal(&invalidStructure)
				if err != nil {
					log.Error().Err(err).Any("user_id", c.id).Msg("Error while marshalling invalid event")
					c.send <- []byte(ERROR_PARSE_MSG)
					continue
				}

				c.send <- invalidStructureBytes
				continue
			}

			// send the message to the room
			var msgTobeSent = EventUserReceivedMessage{
				Event: Event{
					Event: CHAT_GROUP_USER_RECV_MSG,
				},
				Data: Message{
					Id:        ulid.Make().String(),
					UserId:    c.id,
					Name:      c.name,
					Message:   payload.Data.Message,
					CreatedAt: time.Now().UTC(),
				},
			}

			msgBytes, err := json.Marshal(&msgTobeSent)
			if err != nil {
				log.Error().Err(err).Any("user_id", c.id).Msg("Error while marshalling message")
				c.send <- []byte(ERROR_PARSE_MSG)
				continue
			}

			// TODO: save the message to the database
			c.room.broadcast <- msgBytes
		}
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Error().Err(err).Any("user_id", c.id).Msg("Error while writing message")
			break
		}
	}
}
