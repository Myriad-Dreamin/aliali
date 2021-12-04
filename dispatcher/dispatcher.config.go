package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/model"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type ConfigManager struct {
	S suppress.ISuppress
}

func (d *ConfigManager) ReadConfig(configPath string) *ali_notifier.Config {
	f, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		d.S.Suppress(err)
		return nil
	}

	var cfg = new(ali_notifier.Config)
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		d.S.Suppress(err)
	}

	return cfg
}

func (d *Dispatcher) GetConfig() *ali_notifier.Config {
	return d.cfg
}

func (d *Dispatcher) GetDatabase() *gorm.DB {
	return d.db
}

func (d *Dispatcher) chunkSize() int64 {
	if d.cfg == nil {
		return DefaultChunkSize
	}

	if cs, err := strconv.ParseInt(d.cfg.AliDrive.ChunkSize, 10, 64); err != nil {
		d.s.WarnOnce(err)
		return DefaultChunkSize
	} else {
		return cs
	}
}

func (d *Dispatcher) authExpiredX(m *model.AliAuthModel) bool {
	return m == nil || m.ExpiresLocal <= time.Now().Unix()+60
}

func (d *Dispatcher) authExpired() bool {
	return d.authExpiredX(d.auth)
}

func (d *Dispatcher) syncConfig() *ali_notifier.Config {
	return d.cfgMgr.ReadConfig(d.configPath)
}
