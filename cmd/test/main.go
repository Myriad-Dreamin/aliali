package main

import (
  "encoding/json"
  "fmt"
  "github.com/Myriad-Dreamin/aliali/database"
  "github.com/Myriad-Dreamin/aliali/models"
  ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
  "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
  "github.com/Myriad-Dreamin/aliali/pkg/suppress"
  yaml "gopkg.in/yaml.v2"
  "gorm.io/driver/sqlite"
  _ "gorm.io/driver/sqlite"
  "gorm.io/gorm"
  "os"
  "time"
)

func main() {
  ali := ali_drive.NewAli()
  s := suppress.PanicAll{}

  f, err := os.OpenFile("config.yaml", os.O_RDONLY, 0644)
  s.Suppress(err)

  var cfg ali_notifier.Config
  s.Suppress(yaml.NewDecoder(f).Decode(&cfg))

  ali.Headers = append(ali.Headers, [2]string{"origin", "https://aliyundrive.com"})
  ali.Headers = append(ali.Headers, [2]string{"referer", "https://aliyundrive.com"})

  db, err := gorm.Open(sqlite.Open("ali.db"))
  s.Suppress(err)

  s.Suppress(db.AutoMigrate(&models.AliAuthModel{}))

  var c = database.DB{}
  model := &models.AliAuthModel{Key: "primary"}

  if !c.FindAuthModelByKey(db, model) {

    info := ali.Refresh(cfg.AliDrive.RefreshToken)

    b, err := json.Marshal(info)
    s.Suppress(err)

    model.Raw = b
    model.ExpiresLocal = time.Now().Unix() + int64(info.ExpiresIn)
    db.Create(model)
  }

  b := model.Get(s)
  fmt.Println(b.RefreshToken, b.TokenType, b.AccessToken)
}
