package knoxwebhdfs

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/arnobroekhof/knoxwebhdfs/internal/defaults"
	"net"
	"net/http"
	"net/url"
	"time"
)

type AuthTypeKnox string

var (
	AuthTypeBasic AuthTypeKnox = "basic"
	AuthTypeNone  AuthTypeKnox = "none"
)

type Conf struct {
	Addr                  string        `default:"127.0.0.1"`
	Port                  string        `default:"9443"`
	Suffix                string        `default:"gateway"`
	Realm                 string        `default:"default"`
	Scheme                string        `default:"https"`
	BasePath              string        `default:"webhdfs/v1"`
	AuthType              AuthTypeKnox  `default:"none"`
	ConnectionTimeout     time.Duration `default:"30s"`
	ResponseHeaderTimeout time.Duration `default:"30s"`
	DisableKeepAlives     bool          `default:"false"`
	DisableCompression    bool          `default:"true"`
	SSLSkipVerify         bool          `default:"false"`
	MaxIdleConnsPerHost   int
	Username              string
	Password              string
}

type Client struct {
	url       *url.URL
	client    http.Client
	transport *http.Transport
	conf      *Conf
}

func NewClient(conf *Conf) (*Client, error) {
	if conf == nil {
		conf = &Conf{}
	}

	err := defaults.Set(conf)
	if err != nil {
		return nil, fmt.Errorf("unable to set defaults")
	}

	uri := fmt.Sprintf("%s://%s:%s/%s/%s/%s", conf.Scheme, conf.Addr, conf.Port, conf.Suffix, conf.Realm, conf.BasePath)
	url, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to construct url: %s", err)
	}

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: conf.SSLSkipVerify,
		},
		DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, conf.ConnectionTimeout)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		MaxIdleConnsPerHost:   conf.MaxIdleConnsPerHost,
		ResponseHeaderTimeout: conf.ResponseHeaderTimeout,
	}

	h := http.Client{
		Transport: t,
	}

	return &Client{
		url:       url,
		client:    h,
		transport: t,
		conf:      conf,
	}, nil
}
