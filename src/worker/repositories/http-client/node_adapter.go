package http_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sonntuet1997/avalanche-simplified/worker/constants"
	"github.com/sonntuet1997/avalanche-simplified/worker/entities"
	"github.com/sonntuet1997/avalanche-simplified/worker/repositories/http-client/models"
	"github.com/sonntuet1997/avalanche-simplified/worker/utils"
	"gitlab.com/golibs-starter/golib/config"
	"gitlab.com/golibs-starter/golib/web/log"
	"net/http"
	"strings"
)

type NodeRepository struct {
	RestClient    *resty.Client
	AppProperties *config.AppProperties
}

func NewNodeRepository(
	RestClient *resty.Client,
	AppProperties *config.AppProperties,
) *NodeRepository {
	return &NodeRepository{
		RestClient:    RestClient,
		AppProperties: AppProperties,
	}
}

var preferTransactionURLTemplate = "http://%s:%d%sv1/node/prefer-transactions/%d"

func (c *NodeRepository) AskForPreference(ctx context.Context, address string, blockNumber int) (*entities.Transaction, error) {
	var response models.PreferenceResponse
	url := fmt.Sprintf(preferTransactionURLTemplate,
		address, c.AppProperties.Port, c.AppProperties.Path, blockNumber)
	var res *resty.Response
	var err error
	err = utils.Retry(ctx, func() error {
		res, err = c.RestClient.R().SetContext(ctx).Get(url)
		if err != nil {
			log.Debugf("failed to Get Preference (should then retry successfully) with error: %w", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", url, err)
	}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal: %w", err)
	}
	if res.StatusCode() == http.StatusNotFound {
		return nil, constants.ErrRecordNotFound
	}
	if res.StatusCode() != http.StatusOK {
		return nil, errors.New(response.Meta.Message)
	}
	return &response.Data, nil
}

func (c *NodeRepository) CheckHealthAndGetAddress(ctx context.Context, url string) (string, error) {
	res, err := c.RestClient.R().SetContext(ctx).EnableTrace().Get(url)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %w", url, err)
	}
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("not Ok")
	}
	ipAddress := strings.Split(res.Request.TraceInfo().RemoteAddr.String(), ":")[0]
	return ipAddress, nil
}
