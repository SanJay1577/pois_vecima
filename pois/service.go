package pois

import (
	"pois/config"
	alias "pois/pois/alias"
	"pois/pois/ccms"

	"git.eng.vecima.com/cloud/golib/v4/httpservice"
)

// WebService API resource path.
// GetServiceAPIs returns the list of service APIs implemented
// by the metadata module.

func GetServiceAPIs(handler httpservice.HttpServiceHandler, baseResource string) []*httpservice.HttpServiceAPI {
	apis := make([]*httpservice.HttpServiceAPI, 0, 0)
	apis = append(apis, ccms.GetServiceAPIs(handler, baseResource+config.GetConfig().GetString("api.channelpath"))...)
	apis = append(apis, alias.GetServiceAPIs(handler, baseResource+config.GetConfig().GetString("api.aliaspath"))...)
	return apis
}
