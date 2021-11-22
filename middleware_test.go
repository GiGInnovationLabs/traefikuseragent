package traefikuseragent_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	mw "github.com/GiGInnovationLabs/traefikuseragent"
)

func TestBasic(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { called = true })

	instance, err := mw.New(context.TODO(), next, mw.CreateConfig(), "traefikuseragent")
	if err != nil {
		t.Fatalf("Error creating %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

	instance.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("Invalid return code")
	}
	if called != true {
		t.Fatalf("next handler was not called")
	}
}

func TestParse(t *testing.T) {
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	instance, _ := mw.New(context.TODO(), next, mw.CreateConfig(), "traefikuseragent")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Set(mw.UserAgentHeader, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.97 Safari/537.11")

	instance.ServeHTTP(recorder, req)

	assertHeader(t, req, mw.DeviceBotHeader, "false")

	assertHeader(t, req, mw.DeviceMobileHeader, "false")
	assertHeader(t, req, mw.DeviceOsHeader, "Linux")

	assertHeader(t, req, mw.DeviceBrowserHeader, "Chrome")
	assertHeader(t, req, mw.DeviceBrowserVersionHeader, "23.0.1271.97")

	assertHeader(t, req, mw.DeviceEngineHeader, "AppleWebKit")
	assertHeader(t, req, mw.DeviceEngineVersionHeader, "537.11")
}

func TestParseBot(t *testing.T) {
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	instance, _ := mw.New(context.TODO(), next, mw.CreateConfig(), "traefikuseragent")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Set(mw.UserAgentHeader, "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	instance.ServeHTTP(recorder, req)

	assertHeader(t, req, mw.DeviceBotHeader, "true")

	assertHeader(t, req, mw.DeviceMobileHeader, "false")
	assertHeader(t, req, mw.DeviceOsHeader, "")

	assertHeader(t, req, mw.DeviceBrowserHeader, "Googlebot")
	assertHeader(t, req, mw.DeviceBrowserVersionHeader, "2.1")

	assertHeader(t, req, mw.DeviceEngineHeader, "")
	assertHeader(t, req, mw.DeviceEngineVersionHeader, "")
}

func TestParseNoHeader(t *testing.T) {
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	instance, _ := mw.New(context.TODO(), next, mw.CreateConfig(), "traefikuseragent")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

	instance.ServeHTTP(recorder, req)

	assertHeader(t, req, mw.DeviceBotHeader, "false")

	assertHeader(t, req, mw.DeviceMobileHeader, "false")
	assertHeader(t, req, mw.DeviceOsHeader, "")

	assertHeader(t, req, mw.DeviceBrowserHeader, "")
	assertHeader(t, req, mw.DeviceBrowserVersionHeader, "")

	assertHeader(t, req, mw.DeviceEngineHeader, "")
	assertHeader(t, req, mw.DeviceEngineVersionHeader, "")
}

func TestParseBadFormat(t *testing.T) {
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	instance, _ := mw.New(context.TODO(), next, mw.CreateConfig(), "traefikuseragent")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Set(mw.UserAgentHeader, "123asd")

	instance.ServeHTTP(recorder, req)

	assertHeader(t, req, mw.DeviceBotHeader, "false")

	assertHeader(t, req, mw.DeviceMobileHeader, "false")
	assertHeader(t, req, mw.DeviceOsHeader, "")

	assertHeader(t, req, mw.DeviceBrowserHeader, "123asd")
	assertHeader(t, req, mw.DeviceBrowserVersionHeader, "")

	assertHeader(t, req, mw.DeviceEngineHeader, "")
	assertHeader(t, req, mw.DeviceEngineVersionHeader, "")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()
	if req.Header.Get(key) != expected {
		t.Fatalf("invalid value of header [%s] != %s", key, req.Header.Get(key))
	}
}
