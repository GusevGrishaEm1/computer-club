package service

type Config struct {
	FilePath string
}

func NewConfig(filePath string) Config {
	return Config{FilePath: filePath}
}
