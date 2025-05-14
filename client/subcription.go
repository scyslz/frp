package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
)

func setServerAddr(cfg *v1.ClientCommonConfig) {
	// Fetch latest server URL from subscription if provided
	if cfg.SubscriptionUrl != "" {
		if updatedServer, err := fetchServerFromSubscription(cfg.SubscriptionUrl); err == nil {
			log.Infof("Successfully fetched latest server from subscription: %s", updatedServer)
			cfg.ServerAddr, cfg.ServerPort = splitURL(updatedServer)

		} else {
			log.Errorf("Failed to fetch latest server from subscription: %v. Using configured server: %s:%d", err, cfg.ServerAddr, cfg.ServerPort)
		}
	}
}

// Fetches the latest server URL from a subscription address
func fetchServerFromSubscription(subAddr string) (string, error) {
	resp, err := http.Get(subAddr)
	if err != nil {
		return "", fmt.Errorf("failed to perform HTTP GET request to subscription address: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from subscription address: %v", err)
	}
	return strings.TrimSpace(string(body)), nil
}

func splitURL(rawurl string) (string, int) {
	if !strings.Contains(rawurl, "://") {
		rawurl = "http://" + rawurl // Assume http if no scheme is provided
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		// Handle parsing error, perhaps return an error or a default value
		return "", 0
	}

	host := u.Hostname()
	port := u.Port()

	if port == "" {
		switch u.Scheme {
		case "https":
			return host, 443
		case "http":
			return host, 80
		default:
			// Default to 80 for unknown schemes if no port is specified
			return host, 80
		}
	} else {
		p, err := strconv.Atoi(port)
		if err != nil {
			// Handle port conversion error
			return host, 0
		}
		return host, p
	}
}
