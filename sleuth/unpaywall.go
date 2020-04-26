package sleuth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

/*UnpaywallDefaultEndpoint is the standard HTTP endpoint for unpaywall API*/
const UnpaywallDefaultEndpoint = "https://api.unpaywall.org/v2/"

/*UnpaywallResponse represents the response from an unpaywall API call*/
type UnpaywallResponse struct {
	Title   string
	Updated string
	Year    int
	Doi     string
	IsOA    bool     `json:"is_oa"`
	Authors []string `json:"z_authors"`
}

/*UnpaywallClient provides API access to unpaywall API*/
type UnpaywallClient struct {
	email    string
	endpoint string
}

/*GetDefaultUnpaywallClient creates a new UnpaywallClient with default options*/
func GetDefaultUnpaywallClient() (*UnpaywallClient, error) {

	email := os.Getenv("PAPERSLEUTH_EMAIL")

	if email == "" {
		return nil, errors.New("You must set an email via PAPERSLEUTH_EMAIL env var (or create your own client manually)")
	}

	return &UnpaywallClient{email: email, endpoint: UnpaywallDefaultEndpoint}, nil
}

/*GetForDoi returns an UnpaywallResponse for a given DOI string*/
func (client *UnpaywallClient) GetForDoi(doi string) (*UnpaywallResponse, error) {

	req, _ := http.NewRequest("GET", client.endpoint+doi, nil)

	q := url.Values{}
	q.Add("email", client.email)

	req.URL.RawQuery = q.Encode()

	httpClient := http.Client{}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	buff, _ := ioutil.ReadAll(res.Body)

	obj := UnpaywallResponse{}
	json.Unmarshal(buff, &obj)

	defer res.Body.Close()

	return &obj, nil

}
