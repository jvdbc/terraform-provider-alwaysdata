package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-alwaysdata/pkg/api"
)

type Alwaysdata struct {
	httpClient *http.Client
	opts       *AlwaysdataOptions
}

// Read https://pkg.go.dev/net/http#Server
type AlwaysdataOptions struct {
	Endpoint string
	Apikey   string
}

func NewAlwaysdata(c *http.Client, opts *AlwaysdataOptions) *Alwaysdata {
	if opts == nil {
		o := AlwaysdataOptions{Endpoint: "https://api.alwaysdata.com"}
		return &Alwaysdata{httpClient: c, opts: &o}
	}

	return &Alwaysdata{httpClient: c, opts: opts}
}

// https://dev.to/billylkc/parse-json-api-response-in-go-10ng
// https://mholt.github.io/json-to-go/
func (ad *Alwaysdata) Get(id uint) (*api.Database, error) {
	if ad.httpClient == nil {
		return nil, fmt.Errorf("http client must be initialise")
	}

	url := fmt.Sprintf("%v/v1/database/%v", ad.opts.Endpoint, id)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(ad.opts.Apikey, "")

	resp, err := ad.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response %v with %v", resp.Status, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	db := new(api.Database)
	if err := json.Unmarshal(body, db); err != nil {
		return nil, err
	}

	return db, nil
}

func CheckApiKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("AlwaysData API key could not be empty")
	}
	// Others controls
	return nil
}
