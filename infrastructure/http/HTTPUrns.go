package http

import "ci.guzzler.io/guzzler/corcel/core"

var (
	//RequestErrorUrn ...
	RequestErrorUrn = core.ErrorUrn.ForConnector("http").Name("request")

	//ResponseErrorUrn ...
	ResponseErrorUrn = core.ErrorUrn.ForConnector("http").Name("response")

	//RequestURLUrn ...
	RequestURLUrn = core.NewUrn("http").Name("request", "url")

	//RequestBytesUrn ...
	RequestBytesUrn = core.BytesSentCountUrn.ForConnector("http")

	//ResponseBytesUrn ...
	ResponseBytesUrn = core.BytesReceivedCountUrn.ForConnector("http")

	//ResponseStatusUrn ...
	ResponseStatusUrn = core.NewUrn("http").Name("response", "status")

	//RequestHeadersUrn ...
	RequestHeadersUrn = core.NewUrn("http").Name("request", "headers")
)
