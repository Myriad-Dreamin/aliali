package main

import (
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"time"
)

func (w *Worker) chunkSize() int64 {
	if w.cfg == nil {
		return DefaultChunkSize
	}

	if cs, err := strconv.ParseInt(w.cfg.AliDrive.ChunkSize, 10, 64); err != nil {
		w.s.WarnOnce(err)
		return DefaultChunkSize
	} else {
		return cs
	}
}

func (w *Worker) authExpired() bool {
	return w.auth == nil || w.auth.ExpiresLocal <= time.Now().Unix()+60
}

func (w *Worker) syncConfig() *ali_notifier.Config {
	f, err := os.OpenFile(w.configPath, os.O_RDONLY, 0644)
	if err != nil {
		w.s.Suppress(err)
		return nil
	}

	var cfg = new(ali_notifier.Config)
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		w.s.Suppress(err)
	}

	return cfg
}
