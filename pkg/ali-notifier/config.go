package ali_notifier


type AliDriveConfig struct {
  RefreshToken string `yaml:"refresh-token"`
  DriveId      string `yaml:"drive-id"`
  RootPath     string `yaml:"root-path"`
  ChunkSize    string `yaml:"chunk-size"`
}

type Config struct {
  Version  string         `yaml:"version"`
  AliDrive AliDriveConfig `yaml:"ali-drive"`
}
