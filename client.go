package kalshi

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"golang.org/x/time/rate"
)

const (
	APIDemoURL = "https://demo-api.kalshi.co/trade-api/v2/"
	APIProdURL = "https://trading-api.kalshi.com/trade-api/v2/"
)

type Cents int

func (c Cents) String() string {
	dollars := float32(c) / 100
	return fmt.Sprintf("$%.2f", dollars)
}

// Client must be instantiated via NewClient.
type Client struct {
	// BaseURL is one of APIDemoURL or APIProdURL.
	BaseURL string

	// See https://trading-api.readme.io/reference/tiers-and-rate-limits.
	WriteRatelimit *rate.Limiter
	ReadRateLimit  *rate.Limiter

	keyID      string
	privateKey *rsa.PrivateKey
	httpClient *http.Client
}

type CursorResponse struct {
	Cursor string `json:"cursor"`
}

type CursorRequest struct {
	Limit  int    `url:"limit,omitempty"`
	Cursor string `url:"cursor,omitempty"`
}

type request struct {
	CursorRequest
	Method       string
	Endpoint     string
	QueryParams  any
	JSONRequest  any
	JSONResponse any
}

func jsonRequestHeaders(
	ctx context.Context,
	client *http.Client,
	headers http.Header,
	method string, reqURL string,
	jsonReq any, jsonResp any,
) error {
	var reqBodyByt []byte
	var err error
	var body io.Reader
	if jsonReq != nil {
		reqBodyByt, err = json.Marshal(jsonReq)
		if err != nil {
			return err
		}
		body = bytes.NewReader(reqBodyByt)
	}

	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return err
	}
	if headers != nil {
		req.Header = headers
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyByt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reqDump, err := httputil.DumpRequest(req, false)
	if err != nil {
		return err
	}

	respDump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return fmt.Errorf("dump: %w", err)
	}
	dumpErr := fmt.Sprintf("Request\n%s%s\nResponse\n%s%s",
		reqDump,
		reqBodyByt,
		respDump,
		respBodyByt,
	)

	if resp.StatusCode >= 400 {
		return fmt.Errorf(
			"unexpected status: %s\n%s",
			resp.Status,
			dumpErr,
		)
	} else if os.Getenv("KALSHI_HTTP_DEBUG") != "" {
		fmt.Printf("REQUEST DUMP\n%s\n", dumpErr)
	}

	if jsonResp != nil {
		err = json.Unmarshal(respBodyByt, jsonResp)
		if err != nil {
			return fmt.Errorf("unmarshal: %w\n%s", err, dumpErr)
		}
	}
	return nil
}

func (c *Client) request(
	ctx context.Context, r request,
) error {
	u, err := url.Parse(c.BaseURL + r.Endpoint)
	if err != nil {
		return err
	}

	if r.QueryParams != nil {
		v, err := query.Values(r.QueryParams)
		if err != nil {
			return err
		}
		u.RawQuery = v.Encode()
	}

	// Reads: block until rate limit allows (pagination needs many sequential calls).
	// Writes: non-blocking — trades must be fast; fail immediately if rate-limited.
	if r.Method == "GET" {
		if err := c.ReadRateLimit.Wait(ctx); err != nil {
			return fmt.Errorf("read ratelimit: %w", err)
		}
	} else {
		if !c.WriteRatelimit.Allow() {
			return fmt.Errorf("write ratelimit exceeded")
		}
	}

	path := u.Path
	if u.RawQuery != "" {
		path = path + "?" + u.RawQuery
	}
	headers := c.requestHeaders(r.Method, path)
	httpHeaders := make(http.Header)
	for k, v := range headers {
		httpHeaders.Set(k, v)
	}

	return jsonRequestHeaders(
		ctx,
		c.httpClient,
		httpHeaders,
		r.Method,
		u.String(), r.JSONRequest, r.JSONResponse,
	)
}

func (c *Client) requestHeaders(method, path string) map[string]string {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	pathParts := strings.SplitN(path, "?", 2)
	msgString := timestamp + method + pathParts[0]
	signature, err := c.signPSSText(msgString)
	if err != nil {
		// In production code, this should be handled more gracefully
		panic(fmt.Sprintf("failed to sign request: %v", err))
	}

	return map[string]string{
		"KALSHI-ACCESS-KEY":       c.keyID,
		"KALSHI-ACCESS-SIGNATURE": signature,
		"KALSHI-ACCESS-TIMESTAMP": timestamp,
	}
}

// signPSSText signs the given text using RSA-PSS.
func (c *Client) signPSSText(text string) (string, error) {
	hash := sha256.Sum256([]byte(text))

	signature, err := rsa.SignPSS(rand.Reader, c.privateKey, crypto.SHA256, hash[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return "", fmt.Errorf("RSA sign PSS failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// Timestamp represents a POSIX Timestamp in seconds.
type Timestamp time.Time

func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(i), 0))
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(time.Time(t).UTC().Unix()))), nil
}

func basicRateLimit() *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Second), 10)
}

// readRateLimit returns a more generous limiter for read operations.
// Pagination (trending, analytics) requires many sequential calls; 10/s
// burst is too small when the background warmer and user requests share
// the same limiter.
func readRateLimit() *rate.Limiter {
	return rate.NewLimiter(rate.Every(100*time.Millisecond), 30)
}

// NewClient creates a new Kalshi client with RSA key-based authentication.
// The client is ready to use immediately without needing to call Login.
func NewClient(keyID string, privateKey *rsa.PrivateKey, baseURL string) *Client {
	c := &Client{
		httpClient: &http.Client{},
		BaseURL:    baseURL,
		keyID:      keyID,
		privateKey: privateKey,
		// See https://trading-api.readme.io/reference/tiers-and-rate-limits.
		WriteRatelimit: basicRateLimit(),
		ReadRateLimit:  readRateLimit(),
	}

	return c
}

// Deprecated: New creates a client without authentication. Use NewClient instead.
// This function is kept for backward compatibility but will not work with the API.
func New(baseURL string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	return &Client{
		httpClient: &http.Client{
			Jar: jar,
		},
		BaseURL:        baseURL,
		WriteRatelimit: basicRateLimit(),
		ReadRateLimit:  basicRateLimit(),
	}
}

// Time is a time.Time that tolerates additional '"' characters.
// Kalshi API endpoints use both RFC3339 and POSIX
// timestamps.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(bytes.Trim(b, "\"")) == 0 {
		return nil
	}
	err := t.Time.UnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%v: %w", len(b), err)
	}
	return nil
}

// Side is either Yes or No.
type Side string

const (
	Yes Side = "yes"
	No  Side = "no"
)

// SideBool turns a Yes bool into a Side.
func SideBool(yes bool) Side {
	if yes {
		return Yes
	}
	return No
}
