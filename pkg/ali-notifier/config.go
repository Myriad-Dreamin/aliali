package ali_notifier

type AliDriveConfig struct {
	RefreshToken string `yaml:"refresh-token"`
	DriveId      string `yaml:"drive-id"`
	RootPath     string `yaml:"root-path"`
	ChunkSize    string `yaml:"chunk-size"`
}

type BackendConfig struct {
	JwtSecret    string `yaml:"jwt-secret"`
	Account      string `yaml:"account"`
	PasswordHash string `yaml:"password-hash"`
}

type Config struct {
	Version  string         `yaml:"version"`
	AliDrive AliDriveConfig `yaml:"ali-drive"`
	Backend  BackendConfig  `yaml:"backend"`
}
