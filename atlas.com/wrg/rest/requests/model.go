package requests

import (
	"atlas-wrg/rest/response"
	"encoding/json"
	"strconv"
)

type ErrorListDataContainer struct {
	Errors []ErrorData `json:"errors"`
}

type ErrorData struct {
	Status int               `json:"status"`
	Code   string            `json:"code"`
	Title  string            `json:"title"`
	Detail string            `json:"detail"`
	Meta   map[string]string `json:"meta"`
}

type DataContainer[A any] interface {
	Data() DataBody[A]
	DataList() []DataBody[A]
	Included() []interface{}
	Length() int
}

type dataContainer[A any] struct {
	data            response.DataSegment
	included        response.DataSegment
	includedMappers []response.ConditionalMapperProvider
}

type DataBody[A any] struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes A      `json:"attributes"`
}

func (c *dataContainer[A]) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(c.data) == 1 {
		t.Data = c.data[0]
	} else {
		t.Data = c.data
	}
	return json.Marshal(t)
}

func (c *dataContainer[A]) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyDataBody[A]), c.includedMappers...)
	if err != nil {
		return err
	}

	c.data = d
	c.included = i
	return nil
}

func (c dataContainer[A]) Data() DataBody[A] {
	if len(c.data) >= 1 {
		return *c.data[0].(*DataBody[A])
	}
	return DataBody[A]{}
}

func (c dataContainer[A]) DataList() []DataBody[A] {
	var r = make([]DataBody[A], 0)
	for _, x := range c.data {
		r = append(r, *x.(*DataBody[A]))
	}
	return r
}

func (c dataContainer[A]) Included() []interface{} {
	return c.included
}

func (c dataContainer[A]) Length() int {
	return len(c.data)
}

type IncludeFilter[E any] func(i DataBody[E]) bool

func GetInclude[A any, E any](c DataContainer[A], id int) (E, bool) {
	var e E
	for _, x := range c.Included() {
		if val, ok := x.(*DataBody[E]); ok {
			eid, err := strconv.Atoi(val.Id)
			if err == nil && eid == id {
				e = val.Attributes
				return e, true
			}
		}
	}
	return e, false
}

func GetIncluded[A any, E any](c DataContainer[A], filters ...IncludeFilter[E]) []E {
	var e = make([]E, 0)
	for _, x := range c.Included() {
		if val, ok := x.(*DataBody[E]); ok {
			for _, f := range filters {
				ok = f(*val)
			}
			if ok {
				e = append(e, val.Attributes)
			}
		}
	}
	return e
}

func EmptyDataBody[A any]() interface{} {
	return &DataBody[A]{}
}
