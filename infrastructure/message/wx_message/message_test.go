package wx_message

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestXML(t *testing.T) {
	var m WxMessageReply
	m.Content = "text"
	data, _ := xml.Marshal(m)
	fmt.Printf("data: %s", data)
}
