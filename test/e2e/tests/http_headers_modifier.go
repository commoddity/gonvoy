//go:build e2e
// +build e2e

package tests

import (
	"net/http"
	"testing"

	"github.com/commoddity/gonvoy/pkg/suite"
	"github.com/stretchr/testify/require"
)

func init() {
	TestCases = append(TestCases, HttpHeadersModifierTestCase)
}

var HttpHeadersModifierTestCase = suite.TestCase{
	Name:        "HTTPHeadersModifierTest",
	FilterName:  "http_headers_modifier",
	Description: "Running test to simulate HTTP headers modification both Request and Response, while also showing how to use child config for a specific route.",
	Parallel:    true,
	Test: func(t *testing.T, kit *suite.TestSuiteKit) {
		stop := kit.StartEnvoy(t)
		defer stop()

		t.Run("invoke to index route", func(t *testing.T) {
			resp, err := http.Get(kit.GetEnvoyHost())
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, resp.Header.Get("x-header-modified-at"), "parent")
			require.Eventually(t, func() bool {
				return kit.CheckEnvoyLog("request header ---> X-Foo=[\"bar\"]")
			}, kit.WaitDuration, kit.TickDuration, "failed to find log message in envoy log")
		})

		t.Run("invoke to details route, expect to use child config", func(t *testing.T) {
			resp, err := http.Get(kit.GetEnvoyHost() + "/details")
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, resp.Header.Get("x-header-modified-at"), "child")
			require.Eventually(t, func() bool {
				return kit.CheckEnvoyLog("request header ---> X-Boo=[\"far\"]")
			}, kit.WaitDuration, kit.TickDuration, "failed to find log message in envoy log")
		})
	},
}
