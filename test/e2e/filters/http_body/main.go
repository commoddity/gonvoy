package main

import (
	"github.com/commoddity/gonvoy"
)

func init() {
	gonvoy.RunHttpFilter(filterName,
		func() gonvoy.HttpFilter {
			return new(BodyReadFilter)
		},
		gonvoy.ConfigOptions{
			FilterConfig: new(Config),

			DisableStrictBodyAccess: true,
			EnableRequestBodyRead:   true,
			EnableResponseBodyRead:  true,
		},
	)

	gonvoy.RunHttpFilter(filterName,
		func() gonvoy.HttpFilter {
			return new(BodyWriteFilter)
		},
		gonvoy.ConfigOptions{
			FilterConfig: new(Config),

			DisableStrictBodyAccess: true,
			EnableRequestBodyWrite:  true,
			EnableResponseBodyWrite: true,
		},
	)

	gonvoy.RunHttpFilter(filterName,
		func() gonvoy.HttpFilter {
			return new(Echoserver)
		},
		gonvoy.ConfigOptions{
			DisableStrictBodyAccess: true,
			EnableRequestBodyRead:   true,
		},
	)
}

func main() {}
