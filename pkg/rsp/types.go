package rsp

import "github.com/aiechoic/admin/pkg/errs"

type Response struct {
	Success bool      `json:"success"`
	Error   string    `json:"error"`
	Code    errs.Code `json:"code"`
	Data    any       `json:"data"`
}
