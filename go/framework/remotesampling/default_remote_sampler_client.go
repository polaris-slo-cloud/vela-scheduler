package remotesampling

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/go-logr/logr"
	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/client"
)

var (
	_ RemoteSamplerClient = (*DefaultRemoteSamplerClient)(nil)
)

const (
	samplingEndpointsPath = "samples"
)

// Default implementation of RemoteSamplerClient.
type DefaultRemoteSamplerClient struct {
	clusterName          string
	baseURI              string
	samplingStrategyName string
	requestURI           string
	httpClient           *http.Client
	logger               *logr.Logger
}

func NewDefaultRemoteSamplerClient(
	clusterName string,
	baseURI string,
	samplingStrategyName string,
	logger *logr.Logger,
) *DefaultRemoteSamplerClient {
	requestURI, err := url.JoinPath(baseURI, samplingEndpointsPath, samplingStrategyName)
	if err != nil {
		panic(err)
	}

	samplerClient := &DefaultRemoteSamplerClient{
		clusterName:          clusterName,
		baseURI:              baseURI,
		samplingStrategyName: samplingStrategyName,
		requestURI:           requestURI,
		httpClient:           &http.Client{},
		logger:               logger,
	}

	return samplerClient
}

func (sc *DefaultRemoteSamplerClient) BaseURI() string {
	return sc.baseURI
}

func (sc *DefaultRemoteSamplerClient) ClusterName() string {
	return sc.clusterName
}

func (sc *DefaultRemoteSamplerClient) SamplingStrategyName() string {
	return sc.samplingStrategyName
}

func (sc *DefaultRemoteSamplerClient) SampleNodes(ctx context.Context, request *RemoteNodesSamplerRequest) (*RemoteNodesSamplerResponse, *RemoteNodesSamplerError) {
	httpReq, err := sc.createSamplerNodesHttpRequest(ctx, request)
	if err != nil {
		return nil, sc.createSamplerError(err)
	}

	httpResp, err := sc.httpClient.Do(httpReq)
	if err != nil {
		return nil, sc.createSamplerError(err)
	}
	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if httpResp.StatusCode == http.StatusOK {
		if samplerResponse, err := sc.parseSuccessResponseBody(httpResp); err == nil {
			return samplerResponse, nil
		} else {
			return nil, sc.createSamplerError(err)
		}
	} else {
		if samplerError, err := sc.parseErrorResponseBody(httpResp); err == nil {
			return nil, samplerError
		} else {
			return nil, sc.createSamplerError(err)
		}
	}
}

func (sc *DefaultRemoteSamplerClient) createSamplerNodesHttpRequest(ctx context.Context, sampleReq *RemoteNodesSamplerRequest) (*http.Request, error) {
	body, err := json.Marshal(sampleReq)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(body)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", sc.requestURI, bodyBuffer)
	if err != nil {
		return nil, err
	}

	httpReq.Header["Content-Type"] = []string{"application/json"}
	httpReq.Header["Accept"] = []string{"application/json"}

	return httpReq, nil
}

func (sc *DefaultRemoteSamplerClient) createSamplerError(err error) *RemoteNodesSamplerError {
	return &RemoteNodesSamplerError{
		Error: client.NewPolarisErrorDto(err),
	}
}

func (sc *DefaultRemoteSamplerClient) parseSuccessResponseBody(httpResp *http.Response) (*RemoteNodesSamplerResponse, error) {
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	samplerResponse := RemoteNodesSamplerResponse{}
	if err := json.Unmarshal(body, &samplerResponse); err != nil {
		return nil, err
	}

	if samplerResponse.Nodes != nil {
		for _, node := range samplerResponse.Nodes {
			node.ClusterName = sc.clusterName
		}
	}

	return &samplerResponse, nil
}

func (sc *DefaultRemoteSamplerClient) parseErrorResponseBody(httpResp *http.Response) (*RemoteNodesSamplerError, error) {
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	samplerError := RemoteNodesSamplerError{}
	if err := json.Unmarshal(body, &samplerError); err != nil {
		return nil, err
	}

	return &samplerError, nil
}
