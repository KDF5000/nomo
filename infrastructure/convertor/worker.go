package convertor

import (
	"errors"

	"github.com/Jeffail/tunny"
)

var Pool *tunny.Pool

var (
	ErrNoNeedToUpload = errors.New("no need to upload")
	ErrWrongParam     = errors.New("wrong params given")
)

type ConvertOutput struct {
	// convert result if not need to upload to a public cloud storage
	Buf []byte
	// public url if upload to cloud storage successfully
	Url string
	Err error
}

func ConvertHandler(params interface{}) interface{} {
	// doctronOutputDTO := DoctronOutputDTO{}
	var output ConvertOutput
	config, ok := params.(ConvConfig)
	if !ok {
		output.Err = ErrWrongParam
		return output
	}

	conv := NewHtml2ImageConvertor(config)
	data, err := conv.Convert(config.Ctx)
	if err != nil {
		output.Err = err
		return output
	}

	output.Buf = data
	// TODO(kongdefei): upload to cloud storage if needed
	return output
}
