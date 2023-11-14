package vicartt

import "encoding/json"

type Client struct {
	Name      string
	AccessKey string
	conn      chan string
}

func (c *Client) Send(msg interface{}) {
	jsonText, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	c.conn <- string(jsonText)
}
