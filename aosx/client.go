package aosx

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
    "net/http/httputil"
	"time"
)

type Client struct {
	Username   string
	Password   string
	HttpClient *http.Client
}

func NewClient(username, password string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 5*time.Second)
		},
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   60 * time.Second,
	}

	client := &Client{
		Username:   username,
		Password:   password,
		HttpClient: httpClient,
	}

	return client, nil
}

func (c *Client) CreateRestconf(ctx context.Context, path, content string) error {
	reqBody := []byte(content)

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/yang-data+json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) ReadRestconf(ctx context.Context, path string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/yang-data+json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (c *Client) UpdateRestconf(ctx context.Context, path, content string) error {
	reqBody := []byte(content)

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/yang-data+json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteRestconf(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

    req.Header.Set("Accept", "application/yang-data+json")
	req.SetBasicAuth(c.Username, c.Password)

    // Debug: Print the request
	dump, err := httputil.DumpRequestOut(req, true)
    fmt.Printf("REQUEST:\n%s\n", dump)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %v", err)
		}
		bodyString := string(bodyBytes)
		return fmt.Errorf("path is %s, unexpected status code: %d, REQUEST: %s, response body: %s", path, resp.StatusCode, dump, bodyString)
		//return fmt.Errorf("Path is "+path+"unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
