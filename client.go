package sendowl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

// DefaultEndpoint is the default root API endpoint.
const DefaultEndpoint = "https://www.sendowl.com/api/v1/"

var defaultEndpointURL *url.URL

func init() {
	var err error
	defaultEndpointURL, err = url.Parse(DefaultEndpoint)
	if err != nil {
		panic(err)
	}
}

var ErrNotFound = errors.New("not found")

type ResponseNotJSONError struct {
	ContentType string
}

func (e *ResponseNotJSONError) Error() string {
	return fmt.Sprintf("response is not JSON (got Content-Type %q)", e.ContentType)
}

// New creates a new Client which can be used to make requests to Sendowl
// services.
func New(key, secret string) *Client {
	return &Client{
		logger:        log.New(ioutil.Discard, "", log.LstdFlags),
		transportFunc: defaultTransportFunc,
		key:           key,
		secret:        secret,
		endpoint:      defaultEndpointURL,
	}
}

type TransportFunc func(context.Context) http.RoundTripper

func defaultTransportFunc(ctx context.Context) http.RoundTripper {
	return http.DefaultTransport
}

// Client is a type which makes requests to Sendowl.
type Client struct {
	logger        *log.Logger
	transportFunc TransportFunc
	key           string
	secret        string
	endpoint      *url.URL `datastore:"-"`
}

func (c *Client) WithLogger(l *log.Logger) *Client {
	c.logger = l
	return c
}

func (c *Client) WithTransportFunc(f TransportFunc) *Client {
	c.transportFunc = f
	return c
}

func (c *Client) WithEndpoint(e *url.URL) *Client {
	c.endpoint = e
	return c
}

func (c *Client) newRequest(method, refURL string, body io.Reader) (*http.Request, error) {
	ref, err := url.Parse(refURL)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(method, c.endpoint.ResolveReference(ref).String(), body)
	if err != nil {
		return nil, err
	}
	r.SetBasicAuth(c.key, c.secret)
	r.Header.Set("Accept", "application/json")
	return r, nil
}

func (c *Client) do(ctx context.Context, r *http.Request, data interface{}) error {
	c.logger.Printf("sendowl: %s %s (content-type: %q)", r.Method, r.URL, r.Header.Get("Content-Type"))
	rawReq, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return err
	}
	c.logger.Printf("%s", rawReq)
	resp, err := c.transportFunc(ctx).RoundTrip(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.decodeResponse(resp, data)
}

// decodeResponse decodes the response from Sendowl as JSON.
func (c *Client) decodeResponse(resp *http.Response, data interface{}) error {
	body := &bytes.Buffer{}
	r := io.TeeReader(resp.Body, body)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	c.logger.Printf("%s", b)
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return &ResponseNotJSONError{ContentType: ct}
	}
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("sendowl returned non-2xx status %d", resp.StatusCode)
	}
	return json.NewDecoder(body).Decode(data)
}
