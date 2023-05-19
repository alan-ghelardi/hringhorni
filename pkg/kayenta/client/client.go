package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nubank/hringhorni/pkg/kayenta/apis/models"
)

const (

	// Default Kayenta address.
	defaultKayentaAddress = "http://kayenta.h8i-system.svc.cluster.local:8090"

	// Timeout to call the Kayenta API.
	defaultTimeout = 15 * time.Second

	// Kayenta endpoint to create canary analysis.
	canaryEndpoint = "canary"
)

type KayentaClient struct {

	// Kayenta's Base address.
	baseAddress string
}

func (k *KayentaClient) CreateCanaryAnalysis(ctx context.Context, canaryRequest *models.CanaryAdhocExecutionRequest) (*models.CanaryExecutionResponse, error) {
	data, err := json.Marshal(canaryRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling CanaryAnalysisAdhocExecutionRequest object: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", k.baseAddress, canaryEndpoint), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	canaryExecutionResponse := new(models.CanaryExecutionResponse)
	if err := k.doRequest(ctx, request, canaryExecutionResponse); err != nil {
		return nil, err
	}
	return canaryExecutionResponse, nil
}

func (k *KayentaClient) doRequest(ctx context.Context, request *http.Request, object any) error {
	requestCtx, cancelFunc := context.WithTimeout(ctx, defaultTimeout)
	defer cancelFunc()

	request = request.WithContext(requestCtx)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	var body []byte
	if response.Body != nil {
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
	}

	if http.StatusOK != response.StatusCode {
		message := fmt.Sprintf("Kayenta returned an unexpected status code (%d)", response.StatusCode)
		if len(body) != 0 {
			message += ": " + string(body)
		}
		return errors.New(message)
	}

	if err := json.Unmarshal(body, object); err != nil {
		return fmt.Errorf("error unmarshaling response body: %w", err)
	}
	return nil
}

func (k *KayentaClient) GetCanaryAnalysis(ctx context.Context, id string) (*models.CanaryExecutionStatusResponse, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", k.baseAddress, canaryEndpoint, id), io.NopCloser(strings.NewReader("")))
	if err != nil {
		return nil, err
	}

	canaryExecutionStatusResponse := new(models.CanaryExecutionStatusResponse)
	if err := k.doRequest(ctx, request, canaryExecutionStatusResponse); err != nil {
		return nil, err
	}
	return canaryExecutionStatusResponse, nil
}

// New creates a new client to talk to Kayenta.
func New() *KayentaClient {
	return &KayentaClient{
		baseAddress: defaultKayentaAddress,
	}
}
