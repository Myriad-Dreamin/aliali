package dispatcher

import (
	"encoding/json"
	"fmt"
	"github.com/Myriad-Dreamin/aliali/database"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"time"
)

func (d *Dispatcher) makeAliClient(suppress suppress.ISuppress) *ali_drive.Ali {
	a := ali_drive.NewAli(suppress)
	a.Headers = d.httpHeaders
	return a
}

// refreshAuth do expire check, and only invoked in main::Dispatcher context
func (d *Dispatcher) refreshAuth() {
	var c = database.DB{}
	m := &model.AliAuthModel{Key: "primary"}

	if !c.FindAuthModelByKey(d.db, m) || d.authExpiredX(m) {
		d.logger.Printf("refresh access token of aliyunpan")
		info := d.ali.RefreshToken(d.cfg.AliDrive.RefreshToken)

		b, err := json.Marshal(info)
		d.s.Suppress(err)

		m.Raw = b
		m.ExpiresLocal = time.Now().Unix() + int64(info.ExpiresIn)
		c.SaveAuthModel(d.db, m)
	} else {
		d.logger.Printf("using the access token of aliyunpan in the cache")
	}

	b := m.Get(d.s)
	tok := fmt.Sprintf("%s %s", b.TokenType, b.AccessToken)
	if d.authedAli == nil {
		newCli := d.makeAliClient(d.s)
		newCli.SetAccessToken(tok)
		d.authedAli = newCli
	} else {
		d.authedAli.SetAccessToken(tok)
	}
	d.auth = m
}
