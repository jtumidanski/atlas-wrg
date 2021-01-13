package handlers

import (
	"atlas-wrg/attributes"
	"context"
	"net/http"
)

func (c *ChannelServer) MiddlewareValidateChannelServer(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		cs := &attributes.InputChannelServer{}
		err := attributes.FromJSON(cs, r.Body)
		if err != nil {
			c.l.Println("[ERROR] deserializing channel server", err)
			rw.WriteHeader(http.StatusBadRequest)
			attributes.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyChannelServer{}, cs)
		r = r.WithContext(ctx)

		f(rw, r)
	})
}
