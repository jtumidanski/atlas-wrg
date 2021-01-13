package handlers

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// The requested resource was not found
// swagger:response notFoundResponse
type notFoundResponseWrapper struct {
}

// swagger:parameters getWorld getChannelServer
type worldIDParamsWrapper struct {
	// The id of the world for which the operation relates
	// in: path
	// required: true
	// minimum: 0
	WorldId int `json:"worldId"`
}

// swagger:parameters unregister getChannelServer
type channelIDParamsWrapper struct {
	// The id of the channel for which the operation relates
	// in: path
	// required: true
	// minimum: 1
	ChannelId int `json:"channelId"`
}
