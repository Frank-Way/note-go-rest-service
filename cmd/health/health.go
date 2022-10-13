package health

import (
	"flag"
	"fmt"
	"github.com/Frank-Way/note-go-rest-service/internal/server"
	"net/http"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.yaml", "path to config path")
}

func main() {
	flag.Parse()

	config := server.NewConfig(configPath)
	_, err := http.Get(fmt.Sprintf("http://%s:%s/health", config.Listen.BindIP, config.Listen.Port))
	if err != nil {
		os.Exit(1)
	}
}
