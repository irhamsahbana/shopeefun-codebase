package websocket

import "time"

/*
	Constants that can be used in the application
*/

const (
	CHAT_GROUP_USER_LEFT     = "chat-group-user-left"
	CHAT_GROUP_USER_JOIN     = "chat-group-user-join"
	CHAT_GROUP_USER_SEND_MSG = "chat-group-user-send-message"
	CHAT_GROUP_USER_RECV_MSG = "chat-group-user-receive-message"

	VALIDATION_ERROR         = "validation-error"
	VALIDATION_ERROR_MESSAGE = "Your event payload is invalid"
	INVALID_EVENT            = "invalid-event"
	INVALID_EVENT_MESSAGE    = "The event type is invalid"

	ERROR_PARSE_MSG = `
	{
		"event": "parse-error",
		"message": "Error while parsing your event message"
	}
`
	ERROR_SAVE_MSG = `
	{
		"event": "save-message-error",
		"message": "Error while saving your message"
	}
`
)

/*
	*
	Events that can be received from the client
	*
*/

// EventUserSendMessage represents a send message event
// when a user sends a message to the server to be broadcasted
type EventUserSendMessage struct {
	Event
	Data SentMessage `json:"data"`
}

type SentMessage struct {
	Message string `json:"message" validate:"required"`
}

/*
	*
	Events that can be sent to the client
	*
*/

// event represents a message event
type Event struct {
	Event string `json:"event" validate:"required,oneof=chat-group-user-left chat-group-user-join chat-group-user-send-message chat-group-user-receive-message"`
}

type EventInvalid struct {
	Event   string `json:"event"`
	Errors  any    `json:"errors,omitempty"`
	Message string `json:"message"`
}

// EventReceivedMessage represents a received message event
// when a user receives a message from the server
type EventUserReceivedMessage struct {
	Event
	Data Message `json:"data"`
}

type Message struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

/*
	*
	Events that can be sent to the client and received from the client
	*
*/

// EventUserActivity represents a user activity event
// when a user joins or leaves the chat
type EventUserActivity struct {
	Event
	Data UserActivity `json:"data"`
}

type UserActivity struct {
	UserId  string `json:"user_id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
