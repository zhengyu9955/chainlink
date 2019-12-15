package services

import (
	"chainlink/core/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"
)

// Fetcher is the interface encapsulating all functionality needed to retrieve
// a price.
type Fetcher interface {
	Fetch() (float64, error)
}

// httpFetcher retrieves data via HTTP from an external price adapter source.
type httpFetcher struct {
	client      *http.Client
	url         *url.URL
	requestData string
}

func newHTTPFetcher(
	timeout time.Duration,
	requestData,
	urlStr string,
) (Fetcher, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}

	if !gjson.Valid(requestData) {
		return nil, errors.New("unable to parse request data JSON")
	}
	return &httpFetcher{
		client:      &http.Client{Timeout: timeout},
		url:         u,
		requestData: requestData,
	}, err
}

func (p *httpFetcher) Fetch() (float64, error) {
	r, err := p.client.Post(p.url.String(), "application/json", strings.NewReader(p.requestData))
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("unable to fetch price from %s with payload '%s'", p.url.String(), p.requestData))
	}

	defer r.Body.Close()
	target := adapterResponse{}
	if err = json.NewDecoder(r.Body).Decode(&target); err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("unable to decode price from %s", p.url.String()))
	}
	if target.ErrorMessage.Valid {
		return 0, errors.Wrap(errors.New(target.ErrorMessage.String), fmt.Sprintf("price fetcher %s returned error", p.url.String()))
	}
	if r.StatusCode >= 400 {
		return 0, fmt.Errorf("status code: %d, no error message; unable to retrieve price from %s", r.StatusCode, p.url.String())
	}

	result := target.Result()
	if result == nil {
		return 0, errors.Wrap(errors.New("no result returned"), fmt.Sprintf("unable to fetch price from %s", p.url.String()))
	}
	logger.Debugw(
		fmt.Sprintf("fetched price %f from %s", *result, p.url.String()),
		"price", result,
		"url", p.url.String(),
	)
	return *result, nil
}

func (p *httpFetcher) String() string {
	return fmt.Sprintf("http price fetcher: %s", p.url.String())
}

type adapterResponseData struct {
	Result *float64 `json:"result"`
}

// adapterResponse is the HTTP response as defined by the external adapter:
// https://github.com/smartcontractkit/bnc-adapter
type adapterResponse struct {
	Data         adapterResponseData `json:"data"`
	ErrorMessage null.String         `json:"errorMessage"`
}

func (pr adapterResponse) Result() *float64 {
	return pr.Data.Result
}

// medianFetcher fetches from all fetchers, and returns the median value, or
// average if even number of results.
type medianFetcher struct {
	fetchers []Fetcher
}

// newMedianFetcherFromURLs creates a median fetcher that retrieves a price
// from all passed URLs using httpFetcher, and returns the median.
func newMedianFetcherFromURLs(
	timeout time.Duration,
	requestData string,
	priceURLs ...string,
) (Fetcher, error) {
	fetchers := []Fetcher{}
	for _, url := range priceURLs {
		ps, err := newHTTPFetcher(timeout, requestData, url)
		if err != nil {
			return nil, err
		}

		fetchers = append(fetchers, ps)
	}

	medianFetcher, err := newMedianFetcher(fetchers...)
	if err != nil {
		return nil, err
	}

	return medianFetcher, nil
}

func newMedianFetcher(fetchers ...Fetcher) (Fetcher, error) {
	if len(fetchers) == 0 {
		return nil, errors.New("must pass in at least one price fetcher to newMedianFetcher")
	}
	return &medianFetcher{
		fetchers: fetchers,
	}, nil
}

func (m *medianFetcher) Fetch() (float64, error) {
	var err error
	prices := make([]float64, len(m.fetchers))
	fetchErrors := []error{}
	for i, fetcher := range m.fetchers {
		prices[i], err = fetcher.Fetch()
		if err != nil {
			logger.Error(err)
			fetchErrors = append(fetchErrors, err)
		}
	}

	errorRate := float64(len(fetchErrors)) / float64(len(m.fetchers))
	if errorRate >= 0.5 {
		return 0, errors.Wrap(multierr.Combine(fetchErrors...), "majority of fetchers in median failed")
	}

	sort.Float64s(prices)
	k := len(prices) / 2
	if len(prices)%2 == 1 {
		return prices[k], nil
	}
	return (prices[k] + prices[k-1]) / 2, nil
}

func (m *medianFetcher) String() string {
	fetcherDescriptions := make([]string, len(m.fetchers))
	for i, fetcher := range m.fetchers {
		fetcherDescriptions[i] = fmt.Sprintf("%s", fetcher)
	}
	return fmt.Sprintf("median fetcher: %s", strings.Join(fetcherDescriptions, ","))
}