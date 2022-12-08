package Rrpc

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

var Validator StructValidator = &defaultValidator{}

type StructValidator interface {
	// ValidateStruct 结构体验证，如果错误返回对应的错误信息
	ValidateStruct(any) error
	// Engine 返回对应使用的验证器
	Engine() any
}

type defaultValidator struct {
	one      sync.Once
	validate *validator.Validate
}

func (d *defaultValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return d.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return d.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := d.validateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

func (d *defaultValidator) validateStruct(obj any) error {
	d.lazyInit()
	return d.validate.Struct(obj)
}

func (d *defaultValidator) lazyInit() {
	d.one.Do(func() {
		d.validate = validator.New()
	})
}

func (d *defaultValidator) Engine() any {
	d.lazyInit()
	return d.validate
}
