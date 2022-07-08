package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"skuld/encoding"
)

type DecodeErrorFunc func(ctx context.Context, res *http.Response) error

type EncodeRequestFunc func(ctx context.Context, contentType string, in interface{}) (body []byte, err error)

type DecodeResponseFunc func(ctx context.Context, res *http.Response, out interface{}) error

type ClientOption func(*clientOptions)

type clientOptions struct {
	tlsConf      *tls.Config
	timeout      time.Duration
	endpoint     string
	userAgent    string
	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
	transport    http.RoundTripper
}

func WithTlsConfig(cfg *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = cfg
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(o *clientOptions) {
		o.userAgent = userAgent
	}
}

func WithRequestEncoder(encoder EncodeRequestFunc) ClientOption {
	return func(o *clientOptions) {
		o.encoder = encoder
	}
}

func WithResponseDecoder(decoder DecodeResponseFunc) ClientOption {
	return func(o *clientOptions) {
		o.decoder = decoder
	}
}

func WithErrorDecoder(errorDecoder DecodeErrorFunc) ClientOption {
	return func(o *clientOptions) {
		o.errorDecoder = errorDecoder
	}
}

func WithTransport(trans *http.Transport) ClientOption {
	return func(o *clientOptions) {
		o.transport = trans
	}
}

type Client struct {
	opts clientOptions
	cc   *http.Client
}

func NewClient(opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		timeout:      2000 * time.Millisecond,
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		errorDecoder: DefaultErrorDecoder,
		transport:    http.DefaultTransport,
	}

	for _, option := range opts {
		option(&options)
	}

	if options.tlsConf != nil {
		if t, ok := options.transport.(*http.Transport); ok {
			t.TLSClientConfig = options.tlsConf
		}
	}

	return &Client{
		opts: options,
		cc: &http.Client{
			Timeout:   options.timeout,
			Transport: options.transport,
		},
	}, nil
}

func (client *Client) Get(ctx context.Context, path string, reply interface{}) error {
	return client.Invoke(ctx, http.MethodGet, path, "", nil, reply)
}

func (client *Client) Post(ctx context.Context, path string, contentType string, args, reply interface{}) error {
	return client.Invoke(ctx, http.MethodPost, path, contentType, args, reply)
}

// Invoke makes an rpc call procedure for remote service.
func (client *Client) Invoke(ctx context.Context, method, path, contentType string, args interface{}, reply interface{}) error {
	req, err := client.BuildRequest(ctx, method, path, contentType, args)
	if err != nil {
		return err
	}
	return client.invoke(ctx, req, reply)
}

func (client *Client) invoke(ctx context.Context, req *http.Request, reply interface{}) error {
	res, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	err = client.BuildResponse(ctx, res, reply)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) BuildRequest(ctx context.Context, method, path, contentType string, args interface{}) (*http.Request, error) {
	var body io.Reader
	if args != nil {
		if contentType == "" {
			contentType = "application/json"
		}
		data, err := client.opts.encoder(ctx, contentType, args)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	url := fmt.Sprintf("%s%s", client.opts.endpoint, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if client.opts.userAgent != "" {
		req.Header.Set("User-Agent", client.opts.userAgent)
	}
	return req, nil
}

func (client *Client) BuildResponse(ctx context.Context, res *http.Response, reply interface{}) error {
	defer res.Body.Close()
	if err := client.opts.decoder(ctx, res, reply); err != nil {
		return err
	}
	return nil
}

func (client *Client) Do(req *http.Request) (*http.Response, error) {
	return client.do(req.Context(), req)
}

func (client *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := client.cc.Do(req)
	if err != nil {
		return nil, err
	}
	if err = client.opts.errorDecoder(ctx, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func DefaultRequestEncoder(ctx context.Context, contentType string, in interface{}) ([]byte, error) {
	name := ContentSubtype(contentType)
	body, err := encoding.GetCodec(name).Marshal(in)
	if err != nil {
		return nil, err
	}
	return body, err
}

func DefaultResponseDecoder(ctx context.Context, res *http.Response, v interface{}) error {
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return CodecForResponse(res).Unmarshal(data, v)
}

func DefaultErrorDecoder(ctx context.Context, res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err == nil {
		return errors.New(fmt.Sprintf("http status code:%d.", res.StatusCode) + string(data))
	}
	return err
}

func CodecForResponse(r *http.Response) encoding.Codec {
	codec := encoding.GetCodec(ContentSubtype(r.Header.Get("Content-Type")))
	if codec != nil {
		return codec
	}
	return encoding.GetCodec("json")
}
