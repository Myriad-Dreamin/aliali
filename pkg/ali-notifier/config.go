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

type RegistryConfig struct {
	UpstreamSecret string `yaml:"upstream_secret"`
	UpstreamHost   string `yaml:"upstream_host"`
	UpstreamPort   string `yaml:"upstream_port"`

	Name         string `yaml:"name"`
	Upstream     string `yaml:"upstream"`
	Secret       string `yaml:"secret"`
	RegistryHost string `yaml:"registry_host"`
	RegistryPort string `yaml:"registry_port"`
	ServerHost   string `yaml:"server_host"`
	ServerPort   string `yaml:"server_port"`
	Schema       string `yaml:"schema"`
}

type Config struct {
	Version  string                     `yaml:"version"`
	AliDrive AliDriveConfig             `yaml:"ali-drive"`
	Backend  BackendConfig              `yaml:"backend"`
	Servers  map[string]*RegistryConfig `yaml:"servers"`
}
