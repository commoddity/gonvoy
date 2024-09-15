//go:build e2e
// +build e2e

package e2e

import (
	"flag"
	"testing"

	"github.com/commoddity/gonvoy/pkg/suite"
	"github.com/commoddity/gonvoy/test/e2e/tests"
)

func TestE2E(t *testing.T) {
	flag.Parse()

	tSuite := suite.NewTestSuite(suite.TestSuiteOptions{
		EnvoyImageVersion:  suite.DefaultEnvoyImageVersion,
		EnvoyPortStartFrom: 10000,
		AdminPortStartFrom: 8000,
	})

	tSuite.Run(t, tests.TestCases)
}
