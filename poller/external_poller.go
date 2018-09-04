package poller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gbolo/vsummary/common"
)

var (
	// global http client for calls to vsummary server api
	vSummaryClient *http.Client
)

func init() {
	// set sane defaults for vSummaryClient HTTP client
	vSummaryClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   5,
			DisableCompression:    true,
			IdleConnTimeout:       10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}

// ExternalPoller extends Poller with functionality relevant to
// sending results to a vSummary API server over http(s).
type ExternalPoller struct {
	Poller
	stopSignal     chan bool
	vSummaryApiUrl string
}

// NewEmptyExternalPoller returns a empty ExternalPoller
func NewEmptyExternalPoller() *ExternalPoller {
	return &ExternalPoller{
		stopSignal: make(chan bool),
	}
}

// NewExternalPoller returns a ExternalPoller based from a common.Poller
func NewExternalPoller(c common.Poller) (e *ExternalPoller) {
	e = NewEmptyExternalPoller()
	e.Configure(c)
	return
}

// SetEndpoint sets the vSummary API server url unless it's invalid
func (e *ExternalPoller) SetApiUrl(u string) (err error) {
	_, err = url.ParseRequestURI(u)
	if err != nil {
		e.vSummaryApiUrl = u
	}
	return
}

// constructUrl will create the desired vsummary api url
func (e *ExternalPoller) constructUrl(endpoint string) (urlEndpont string, err error) {
	if e.vSummaryApiUrl != "" && endpoint != "" {
		urlEndpont = fmt.Sprintf("%s%s", e.vSummaryApiUrl, endpoint)
	} else {
		err = fmt.Errorf("vSummaryApiUrl or endpoint is empty")
		return
	}
	_, err = url.ParseRequestURI(urlEndpont)
	return
}

// sendResult does an http post request to the vsummary api server to proccess the poll result
func (e *ExternalPoller) sendResult(endpoint string, jsonBody []byte) (err error) {
	// determine url
	url, err := e.constructUrl(endpoint)
	if err != nil {
		return
	}

	// send request
	log.Debugf("sending results to: %s", url)
	res, err := vSummaryClient.Post(url, "application/json", bytes.NewReader(jsonBody))

	// this means the vsummary server api is unreachable
	if err != nil {
		log.Errorf("vsummary api is unreachable: %s error %s", url, err)
		return
	}

	// we only accept 202 as success
	if res.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("recieved %d response code from %", res.StatusCode, url)
		return
	}

	// To ensure KeepAlive:
	// Read until Response is complete (i.e. ioutil.ReadAll(rep.Body))
	// Call Body.Close()
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()

	log.Infof("api call successful: %d %s", res.StatusCode, url)
	return
}
