package domain

type Config struct {
	LogLevel        string `env:"LOG_LEVEL" envDefault:"debug"`
	EnableSwagger   bool   `env:"ENABLE_SWAGGER" envDefault:"true"`
	EnableProfiling bool   `env:"ENABLE_PROFILING" envDefault:"true"`
	HttpPort        int    `env:"HTTP_PORT" envDefault:"8000"`
	HidPath         string `env:"HID_PATH" envDefault:"/dev/hidg0"`
}
