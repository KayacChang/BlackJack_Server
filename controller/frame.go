package controller

import "gitlab.com/ti-backend/go-modules/frame"

type Frame struct {
	frame.Frame
	Token     string `json:"token,omitempty"`
	GameToken string `json:"game_token,omitempty"`
	GameID    string `json:"game_id,omitempty"`
}
