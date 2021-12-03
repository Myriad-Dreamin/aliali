package main

import (
	"encoding/json"
	"fmt"
	"github.com/Myriad-Dreamin/aliali/database"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"time"
)

func (w *Worker) makeAliClient() *ali_drive.Ali {
	a := ali_drive.NewAli()
	a.Headers = w.httpHeaders
	return a
}

// refreshAuth do expire check, and only invoked in main::Worker context
func (w *Worker) refreshAuth() {
	var c = database.DB{}
	m := &model.AliAuthModel{Key: "primary"}

	if !c.FindAuthModelByKey(w.db, m) || w.authExpired() {
		info := w.ali.RefreshToken(w.cfg.AliDrive.RefreshToken)

		b, err := json.Marshal(info)
		w.s.Suppress(err)

		m.Raw = b
		m.ExpiresLocal = time.Now().Unix() + int64(info.ExpiresIn)
		c.SaveAuthModel(w.db, m)
	}

	w.auth = m

	newCli := w.makeAliClient()
	b := m.Get(w.s)
	newCli.SetAccessToken(fmt.Sprintf("%s %s", b.TokenType, b.AccessToken))
	w.authedAli = newCli
}
