// Package traefikuseragent a demo plugin.
package traefikuseragent

import (
	"context"
	"net/http"
	"strconv"

	"github.com/mssola/user_agent"
)

const (
	// UserAgentHeader header.
	UserAgentHeader = "User-Agent"

	// DeviceBotHeader header.
	DeviceBotHeader = "X-Device-Bot"

	// DeviceMobileHeader header.
	DeviceMobileHeader = "X-Device-Mobile"
	// DeviceOsHeader header.
	DeviceOsHeader = "X-Device-Os"

	// DeviceBrowserHeader header.
	DeviceBrowserHeader = "X-Device-Browser"
	// DeviceBrowserVersionHeader header.
	DeviceBrowserVersionHeader = "X-Device-Browser-Version"

	// DeviceEngineHeader header.
	DeviceEngineHeader = "X-Device-Engine"
	// DeviceEngineVersionHeader header.
	DeviceEngineVersionHeader = "X-Device-Engine-Version"
)

// Config the plugin configuration.
type Config struct{}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// TraefikUserAgent a TraefikUserAgent plugin.
type TraefikUserAgent struct {
	next http.Handler
	name string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &TraefikUserAgent{
		next: next,
		name: name,
	}, nil
}

func (mw *TraefikUserAgent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ua := user_agent.New(req.Header.Get(UserAgentHeader))

	req.Header.Set(DeviceBotHeader, strconv.FormatBool(ua.Bot()))

	req.Header.Set(DeviceMobileHeader, strconv.FormatBool(ua.Mobile()))
	req.Header.Set(DeviceOsHeader, ua.OSInfo().Name)

	name, version := ua.Browser()
	req.Header.Set(DeviceBrowserHeader, name)
	req.Header.Set(DeviceBrowserVersionHeader, version)

	name, version = ua.Engine()
	req.Header.Set(DeviceEngineHeader, name)
	req.Header.Set(DeviceEngineVersionHeader, version)

	mw.next.ServeHTTP(rw, req)
}
