package files

import (
	"github.com/pedromspeixoto/users-api/internal/config"
	"github.com/pedromspeixoto/users-api/internal/pkg/logger"
	"go.uber.org/fx"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func ProvideFileServingClient() fx.Option {
	return fx.Provide(NewFileServingClient)
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type fileServingDeps struct {
	fx.In

	Config *config.Config
	Logger *logger.LoggingClient
}

type FileServingClient struct {
	fileServingUrl string
	httpClient
	logger.Logger
}

type clientOption func(c *FileServingClient)

func WithHttpClient(httpClient httpClient) clientOption {
	return func(c *FileServingClient) {
		c.httpClient = httpClient
	}
}

func NewFileServingClient(deps fileServingDeps, opts ...clientOption) (*FileServingClient, error) {
	var err error

	// check url
	_, err = url.Parse(deps.Config.FileServingUrl)
	if err != nil {
		return nil, errors.New("invalid file serving client url provided")
	}

	client := &FileServingClient{
		fileServingUrl: deps.Config.FileServingUrl,
		httpClient:     http.DefaultClient,
		Logger:         deps.Logger.GetLogger(),
	}
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

type RequestOptions struct {
	Logger log.Logger
}

func (c *FileServingClient) GetRandomFile(opts ...RequestOptions) ([]byte, string, error) {

	// fetch the file from the file-serving service
	fileURL := c.fileServingUrl
	response, err := http.Get(fileURL)
	if err != nil {
		c.Logger.Errorf("failed to fetch file from %s. Error: %v", fileURL, err)
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.Logger.Errorf("failed to fetch file from %s. Error: %v", fileURL, err)
		return nil, "", err
	}

	// read the file content from the response body
	fileContent, err := io.ReadAll(response.Body)
	if err != nil {
		c.Logger.Errorf("failed to read file from response body. Error: %v", err)
		return nil, "", err
	}

	// get file type from response header
	fileType := response.Header.Get("Content-Type")

	return fileContent, fileType, nil
}
