package storage

type Config struct {
	Default string      `yaml:"default" json:"default"`
	Local   LocalConfig `yaml:"local" json:"local"`
	Minio   MinioConfig `yaml:"minio" json:"minio"`
	S3      S3Config    `yaml:"s3" json:"s3"`
	DefaultBaseURL string `yaml:"default_base_url" json:"default_base_url"`
}

type LocalConfig struct {
	UploadFolder string `yaml:"upload_folder" json:"upload_folder"`
	BaseURL      string `yaml:"base_url" json:"base_url"`
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint" json:"endpoint"`
	AccessKey string `yaml:"access_key" json:"access_key"`
	SecretKey string `yaml:"secret_key" json:"secret_key"`
	Bucket    string `yaml:"bucket" json:"bucket"`
	Secure    bool   `yaml:"secure" json:"secure"`
	Region    string `yaml:"region" json:"region"`
	BaseURL   string `yaml:"base_url" json:"base_url"`
}

type S3Config struct {
	Endpoint  string `yaml:"endpoint" json:"endpoint"`
	AccessKey string `yaml:"access_key" json:"access_key"`
	SecretKey string `yaml:"secret_key" json:"secret_key"`
	Bucket    string `yaml:"bucket" json:"bucket"`
	Region    string `yaml:"region" json:"region"`
	PathStyle bool   `yaml:"path_style" json:"path_style"`
	BaseURL   string `yaml:"base_url" json:"base_url"`
}
