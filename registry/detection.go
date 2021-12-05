package registry

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"time"
)

type ContainerStatus struct {
	State struct {
		Status     string    `json:"Status"`
		Running    bool      `json:"Running"`
		Paused     bool      `json:"Paused"`
		Restarting bool      `json:"Restarting"`
		OOMKilled  bool      `json:"OOMKilled"`
		Dead       bool      `json:"Dead"`
		Pid        int       `json:"Pid"`
		ExitCode   int       `json:"ExitCode"`
		Error      string    `json:"Error"`
		StartedAt  time.Time `json:"StartedAt"`
		FinishedAt time.Time `json:"FinishedAt"`
	} `json:"State"`
}

func (c *Client) getRecorderStatus() []ContainerStatus {
	var b = bytes.NewBuffer(nil)
	var cmd = exec.Command("docker", "inspect", "bilibili-recorder")
	cmd.Stdout = b
	_ = cmd.Run()

	if b.Len() != 0 {
		var res []ContainerStatus
		_ = json.Unmarshal(b.Bytes(), &res)
		return res
	}

	return nil
}

func (c *Client) DetectRecorderStatus() bool {
	var x = c.getRecorderStatus()
	if len(x) == 0 {
		return false
	}
	if x[0].State.Running == false {
		return false
	}

	return true
}
