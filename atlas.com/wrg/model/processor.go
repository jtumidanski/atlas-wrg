package model

import "errors"

type Operator[M any] func(M) error

type Provider[M any] func() (M, error)

type SliceOperator[M any] func([]M) error

type SliceProvider[M any] func() ([]M, error)

type PreciselyOneFilter[M any] func([]M) (M, error)

func ExecuteForEach[M any](f Operator[M]) SliceOperator[M] {
	return func(models []M) error {
		for _, m := range models {
			err := f(m)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

type Filter[M any] func(M) bool

func FilteredProvider[M any](provider SliceProvider[M], filters ...Filter[M]) SliceProvider[M] {
	models, err := provider()
	if err != nil {
		return ErrorSliceProvider[M](err)
	}

	var results []M
	for _, m := range models {
		good := true
		for _, f := range filters {
			if !f(m) {
				good = false
				break
			}
		}
		if good {
			results = append(results, m)
		}
	}
	return FixedSliceProvider(results)
}

func FixedProvider[M any](model M) Provider[M] {
	return func() (M, error) {
		return model, nil
	}
}

func FixedSliceProvider[M any](models []M) SliceProvider[M] {
	return func() ([]M, error) {
		return models, nil
	}
}

func ErrorProvider[M any](err error) Provider[M] {
	return func() (M, error) {
		var m M
		return m, err
	}
}

func ErrorSliceProvider[M any](err error) SliceProvider[M] {
	return func() ([]M, error) {
		return nil, err
	}
}

func SliceProviderToProviderAdapter[M any](provider SliceProvider[M], preciselyOneFilter PreciselyOneFilter[M]) Provider[M] {
	return func() (M, error) {
		ps, err := provider()
		if err != nil {
			var result M
			return result, err
		}
		return preciselyOneFilter(ps)
	}
}

func IfPresent[M any](provider Provider[M], operator Operator[M]) {
	model, err := provider()
	if err != nil {
		return
	}
	_ = operator(model)
}

func For[M any](provider SliceProvider[M], operator SliceOperator[M]) {
	models, err := provider()
	if err != nil {
		return
	}
	_ = operator(models)
}

func ForEach[M any](provider SliceProvider[M], operator Operator[M]) {
	For(provider, ExecuteForEach(operator))
}

type Transformer[M any, N any] func(M) (N, error)

func Map[M any, N any](provider Provider[M], transformer Transformer[M, N]) Provider[N] {
	m, err := provider()
	if err != nil {
		return ErrorProvider[N](err)
	}
	n, err := transformer(m)
	if err != nil {
		return ErrorProvider[N](err)
	}
	return FixedProvider(n)
}

func SliceMap[M any, N any](provider SliceProvider[M], transformer Transformer[M, N]) SliceProvider[N] {
	models, err := provider()
	if err != nil {
		return ErrorSliceProvider[N](err)
	}
	var results = make([]N, 0)
	for _, m := range models {
		var n N
		n, err = transformer(m)
		if err != nil {
			return ErrorSliceProvider[N](err)
		}
		results = append(results, n)
	}
	return FixedSliceProvider(results)
}

func First[M any](provider SliceProvider[M], filters ...Filter[M]) (M, error) {
	var r M
	ms, err := provider()
	if err != nil {
		return r, err
	}

	if len(filters) == 0 {
		return ms[0], nil
	}

	for _, m := range ms {
		ok := true
		for _, filter := range filters {
			if !filter(m) {
				ok = false
			}
		}
		if ok {
			return m, nil
		}
	}
	return r, errors.New("no result found")
}
