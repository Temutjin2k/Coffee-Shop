package models

import "errors"

var (
	ErrOrderClosed = errors.New("the order is already closed")
)

type Error struct {
	Code    int    `json:"StatusCode"`
	Message string `json:"ErrorMessage"`
}
