package requests

import (
	"atlas-wrg/rest/response"
)

const (
	BaseRequest string = "http://atlas-nginx:80"
)

type configuration struct {
	retries int
	mappers []response.ConditionalMapperProvider
}

type Configurator func(c *configuration)

func SetRetries(amount int) Configurator {
	return func(c *configuration) {
		c.retries = amount
	}
}

func AddMappers(mappers []response.ConditionalMapperProvider) Configurator {
	return func(c *configuration) {
		c.mappers = append(c.mappers, mappers...)
	}
}

func AddMapper(mapper response.ConditionalMapperProvider) Configurator {
	return func(c *configuration) {
		c.mappers = append(c.mappers, mapper)
	}
}
