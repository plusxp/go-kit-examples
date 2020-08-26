package addservice

import (
	"context"
	"log"
	"reflect"
)

func NewServiceWrapper(svc interface{}) *SvcWrapper {
	hd := make(map[string]SvcHandler)
	v := reflect.ValueOf(svc)
	t := reflect.TypeOf(svc)

	for i := 0; i < v.NumMethod(); i++ {
		h, ok := v.Method(i).Interface().(func(context.Context, interface{}) (interface{}, error))
		if ok {
			log.Printf("NewServiceWrapper: method [%s] registed", t.Method(i).Name)
			hd[t.Method(i).Name] = h
		}
	}
	return &SvcWrapper{hd: hd}
}

type SvcWrapper struct {
	hd map[string]SvcHandler
}

type SvcHandler func(context.Context, interface{}) (interface{}, error)

// It may return nil if method not found.
func (s *SvcWrapper) GetHandler(method string) SvcHandler {
	return s.hd[method]
}

func (s *SvcWrapper) Range(f func(method string, h SvcHandler)) {
	for s, h := range s.hd {
		f(s, h)
	}
}
