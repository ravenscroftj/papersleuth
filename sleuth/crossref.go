package sleuth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// TODO: Support rate limiting  as per https://github.com/CrossRef/rest-api-doc#rate-limits

const DefaultCrossRefEndpoint string = "https://api.crossref.org"
const DefaultHomepage string = "https://www.brainsteam.co.uk"
const UserAgentFormatString string = "PaperSleuth/1.0 (%s; mailto:%s)"

type CrossrefClient struct {
	Endpoint string
	Email    string
	Homepage string
}

type crossrefMessage struct {
	*CrossrefWork
}

/*crossRefMessage envelopes results from API*/
type crossRefEnvelope struct {
	Status         string `json:"status"`
	MessageType    string `json:"message-type"`
	MessageVersion string `json:"message-version"`
	Message        crossrefMessage
}

func (e *crossRefEnvelope) UnmarshalJSON(data []byte) error {

	type envtype struct {
		Status         string `json:"status"`
		MessageType    string `json:"message-type"`
		MessageVersion string `json:"message-version"`
	}

	tmp := envtype{}
	json.Unmarshal(data, &tmp)

	e.Status = tmp.Status
	e.MessageType = tmp.MessageType
	e.MessageVersion = tmp.MessageVersion

	switch e.MessageType {

	case "work":
		type WorkEnvelope struct {
			Message CrossrefWork
		}

		val := WorkEnvelope{}
		json.Unmarshal(data, &val)

		e.Message.CrossrefWork = &val.Message

	}

	return nil
}

type CrossrefWork struct {
	Title       []string
	Abstract    string
	DOI         string
	articleType string `json:"type"`
}

func GetDefaultCrossrefClient() (*CrossrefClient, error) {
	email := os.Getenv("PAPERSLEUTH_EMAIL")

	if email == "" {
		return nil, errors.New("You must set an email via PAPERSLEUTH_EMAIL env var (or create your own client manually)")
	}

	return &CrossrefClient{Email: email, Endpoint: DefaultCrossRefEndpoint, Homepage: DefaultHomepage}, nil
}

func (client *CrossrefClient) GetWorkByDOI(doi string) (*CrossrefWork, error) {

	httpClient := http.Client{}

	req, err := http.NewRequest("GET", client.Endpoint+"/works/"+doi, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", fmt.Sprintf(UserAgentFormatString, client.Email, client.Homepage))

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	response := crossRefEnvelope{}

	//try to read response
	buffer, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buffer, &response); err != nil {
		return nil, err
	}

	if response.Status == "error" {
		return nil, errors.New(response.Status + " " + response.MessageType)
	}

	return response.Message.CrossrefWork, nil

}
