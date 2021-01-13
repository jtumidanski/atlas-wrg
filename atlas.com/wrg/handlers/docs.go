package handlers

// No content is returned by this API endpoint
// swagger:response noContentResponse
//goland:noinspection GoUnusedType
type noContentResponseWrapper struct {
}

// The requested resource was not found
// swagger:response notFoundResponse
//goland:noinspection GoUnusedType
type notFoundResponseWrapper struct {
}

// swagger:parameters getWorld getChannelServer
//goland:noinspection GoUnusedType
type worldIDParamsWrapper struct {
	// The id of the world for which the operation relates
	// in: path
	// required: true
	// minimum: 0
	WorldId int `json:"worldId"`
}

// swagger:parameters unregister getChannelServer
//goland:noinspection GoUnusedType
type channelIDParamsWrapper struct {
	// The id of the channel for which the operation relates
	// in: path
	// required: true
	// minimum: 1
	ChannelId int `json:"channelId"`
}
