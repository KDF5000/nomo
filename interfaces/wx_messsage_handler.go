package interfaces

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KDF5000/nomo/application"
	"github.com/KDF5000/nomo/infrastructure/message/wx_message"
	"github.com/KDF5000/pkg/log"
)

type wxMessageHandler struct {
	messageHandleApp application.MessageHandleApp
}

func NewWXMessageHandler(app application.MessageHandleApp) *wxMessageHandler {
	return &wxMessageHandler{messageHandleApp: app}
}

func (h *wxMessageHandler) UrlVerification(c *gin.Context) {
	var event wx_message.WechatVerifyParam
	if err := c.ShouldBind(&event); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.messageHandleApp.VerifyURL(context.TODO(), &event)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, event.Echostr)
}

func (h *wxMessageHandler) HandleMessage(c *gin.Context) {
	if c.Request == nil || c.Request.Body == nil {
		c.String(http.StatusBadRequest, "invalid body")
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Infof("data ==> %+v", string(data))
	var message wx_message.WxMessage
	if err := xml.Unmarshal(data, &message); err != nil {
		log.Warn(fmt.Sprintf("failed to unmarshal message, %v", err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if message.MsgType != "text" {
		log.Warnf("message type %s not supported", message.MsgType)
		c.String(http.StatusBadRequest, "message type not supported")
		return
	}

	reply := wx_message.WxMessageReply{
		ToUserName:   message.FromUserName,
		FromUserName: message.ToUserName,
		CreateTime:   uint64(time.Now().Unix()),
		MsgType:      "text",
		Content:      message.Content,
	}

	c.XML(http.StatusOK, &reply)
}
