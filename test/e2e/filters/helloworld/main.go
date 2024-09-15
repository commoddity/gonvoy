package main

import "github.com/commoddity/gonvoy"

func init() {
	gonvoy.RunHttpFilter(
		helloWorldFilterName,
		func() gonvoy.HttpFilter {
			return new(Filter)
		},
		gonvoy.ConfigOptions{},
	)
}

const helloWorldFilterName = "helloworld"

type Filter struct{}

func (Filter) OnBegin(c gonvoy.RuntimeContext, ctrl gonvoy.HttpFilterController) error {
	c.Log().Info("Hello World from the helloworld HTTP filter")
	return nil
}

func (Filter) OnComplete(c gonvoy.Context) error {
	return nil
}

func main() {}
