package common

type ErrCode uint16

const (
	ErrSucc ErrCode = iota
	ErrInvalidParam
	ErrInternalError
)

type APIResonse struct {
	Code    ErrCode        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
