package registry

import (
	"fmt"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/go-resty/resty/v2"
	"time"
)

type Client struct {
	cli    *resty.Client
	Ident  string
	Config *ali_notifier.RegistryConfig
}

func NewClient(ident string, cfg *ali_notifier.RegistryConfig) *Client {
	return &Client{
		cli:    resty.New(),
		Ident:  ident,
		Config: cfg,
	}
}

func (c *Client) Run() error {
	tick := time.NewTicker(time.Second * 15)
	c.heartbeat()
	for {
		select {
		case <-tick.C:
			c.heartbeat()
		}
	}
}

func (c *Client) heartbeat() {
	var req = PostRegisterRequest{
		Name:   c.Config.Name,
		Secret: c.Config.UpstreamSecret,
		Port:   c.Config.ServerPort,
		Server: c.Ident,
	}
	req.RecorderStatus = 2

	r := c.cli.R().
		SetBody(req)
	resp, err := r.Post(fmt.Sprintf("%s://%s:%s/registry/v1/register", c.Config.Schema, c.Config.UpstreamHost, c.Config.UpstreamPort))
	if err != nil {
		fmt.Printf("heartbeat error: %s\n", err.Error())
		return
	}
	if resp.StatusCode() != 200 {
		fmt.Printf("heartbeat error: %s %s\n", resp.Status(), string(resp.Body()))
	}
	_ = resp
}
