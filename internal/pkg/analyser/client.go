package analyser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/mailru/easyjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ analyse.Analyser = &HTTPClient{}

type HTTPClient struct {
	baseURL  string
	timeout  time.Duration
	log      *zap.Logger
	inputDir string
}

type Config struct {
	BaseURL  string        `yaml:"base_url"`
	Timeout  time.Duration `json:"timeout"`
	InputDir string        `yaml:"input_dir"`
}

type In struct {
	fx.In

	Cfg    Config
	Logger *zap.Logger
}

func New(in In) *HTTPClient {
	return &HTTPClient{
		baseURL:  in.Cfg.BaseURL,
		timeout:  in.Cfg.Timeout,
		log:      in.Logger,
		inputDir: in.Cfg.InputDir,
	}
}

func (c *HTTPClient) Run(ctx context.Context, filename string) (*model.InternalAnalyseResult, error) {
	if c.baseURL == "" {
		return nil, errors.New("base URL is not configured")
	}

	payload := model.InternalAnalyseRequest{
		InputFile: filepath.Join(c.inputDir, filename),
	}

	data, err := easyjson.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: c.timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected response status: " + resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	c.log.Debug("raw body", zap.String("raw body", string(b)))

	var result model.InternalAnalyseResult
	if err := easyjson.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}
