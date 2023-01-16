package handlers

import (
	"fmt"
	"github.com/pin/tftp/v3"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpProxyGetHandler struct {
	baseUrl string
	client  *http.Client
}

func NewHttpProxyGetHandler(baseUrl string, timeOut time.Duration) (*HttpProxyGetHandler, error) {
	return &HttpProxyGetHandler{
		baseUrl: baseUrl,
		client:  &http.Client{Timeout: timeOut},
	}, nil
}

func (hpgh *HttpProxyGetHandler) Handler(filename string, rf io.ReaderFrom) error {
	log.WithField("filename", filename).
		Infof("RRQ received for %s", filename)
	raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()

	proxyUrl, err := url.JoinPath(hpgh.baseUrl, filename)
	if err != nil {
		log.WithField("filename", filename).
			WithError(err).
			Warn("Failed to create proxy URL from filename %s", filename)
		return fmt.Errorf("file not found")
	}

	log.WithField("filename", filename).
		Debugf("Proxying request to %s", proxyUrl)

	req, err := http.NewRequest("GET", proxyUrl, nil)
	if err != nil {
		log.WithField("filename", filename).
			WithError(err).Error("Failed to create request")
		return fmt.Errorf("proxy error")
	}
	req.Header.Add("X-TFTP-CLIENT", raddr.IP.String())

	response, err := hpgh.client.Do(req)
	if err != nil {
		log.WithField("filename", filename).
			WithError(err).Error("Failed to execute request")
		return fmt.Errorf("proxy error")
	}

	if response.StatusCode == http.StatusNotFound {
		// This could be perfectly reasonable
		log.WithField("filename", filename).
			Debugf("File not found on proxy url")
		return fmt.Errorf("file not found")
	}

	if response.StatusCode != http.StatusOK {
		log.WithField("filename", filename).
			Warnf("Response \"%d %s\" received, failing request", response.StatusCode, response.Status)
		return fmt.Errorf("proxy error")
	}

	if response.ContentLength >= 0 {
		rf.(tftp.OutgoingTransfer).SetSize(response.ContentLength)
	}

	n, err := rf.ReadFrom(response.Body)
	if err != nil {
		log.WithField("filename", filename).
			WithError(err).
			Error("Failed to read complete body")
		return fmt.Errorf("write error")
	}

	log.WithField("filename", filename).
		Infof("RRQ %s completed (%d bytes)", filename, n)
	return nil
}
