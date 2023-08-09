package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vjerci/reverse-proxy/internal/block"
	"github.com/vjerci/reverse-proxy/internal/config"
	customlog "github.com/vjerci/reverse-proxy/internal/log"
	"github.com/vjerci/reverse-proxy/internal/mask"
	"github.com/vjerci/reverse-proxy/internal/proxy"
	"github.com/vjerci/reverse-proxy/internal/server"
)

var ErrGuardCreation = errors.New("failed to instantiate blocking guards")

func Build() (http.HandlerFunc, error) {
	configData, err := config.Load(os.Getenv("CONFIG_FILE"))
	if err != nil {
		return nil, err
	}

	log.Printf("proxy forwarding to %s://%s", configData.ForwardScheme, configData.ForwardHost)

	guard, err := block.GuardsFromInterface(configData.Block, &block.InterfaceGuardDecoder{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGuardCreation, err)
	}

	inspector := mask.NewJSONInspector(mask.NewJSONMask(), mask.NewPIIClassifier(mask.NewDefaultPIIPatterns()))
	proxy := proxy.NewProxy(&http.Client{
		Timeout: time.Duration(2 * time.Second),
	})

	responseWriterFactory := &customlog.ResponseWriterFactoryInstance{
		Logger: log.Default(),
	}

	return http.HandlerFunc(server.Handle(inspector, responseWriterFactory, guard, proxy, configData.ForwardHost, configData.ForwardScheme)), nil
}
