package inject

import (
	//"fmt"
	"reflect"
	"strconv"
)

// 函数调用
type Invoke interface {
	// 调用函数，以reflect.Value方式返回参数及错误信息
	Invoke(interface{}) ([]reflect.Value, error)
}

// 赋值
type Assign interface {
	// 结构体字段赋值
	AssignField(interface{}) error
}

// 按照index进行Type映射
type TypeIndexMapper interface {
	MapIndex(interface{}, int) TypeIndexMapper

	MapIndexTo(interface{}, interface{}, int) TypeIndexMapper

	SetIndex(reflect.Type, reflect.Value, int) TypeIndexMapper

	GetIndex(reflect.Type, int) reflect.Value
}

// 按照tag进行Type映射
type TypeTagMapper interface {
	MapTag(interface{}, string) TypeTagMapper

	MapTagTo(interface{}, interface{}, string) TypeTagMapper

	SetTag(reflect.Type, reflect.Value, string) TypeTagMapper

	GetTag(reflect.Type, string) reflect.Value
}

type TypeMapper interface {
	TypeIndexMapper
	//TypeTagMapper
}

type Injector interface {
	Invoke
	Assign
	TypeMapper
}

// 获得指向interface指针的类型
func InterfaceOfPtr(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic("Inject.Interface: Expect a pointer to an interface.")
	}

	return t
}

type injector struct {
	values map[reflect.Type]map[string]reflect.Value
	parent Injector
}

func (i *injector) MapIndex(value interface{}, index int) TypeIndexMapper {
	t := reflect.TypeOf(value)

	m, ok := i.values[t]

	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	tag := strconv.Itoa(index)
	m[tag] = reflect.ValueOf(value)

	return i
}

func (i *injector) MapIndexTo(value interface{}, interfacePtr interface{}, index int) TypeIndexMapper {
	t := InterfaceOfPtr(interfacePtr)

	m, ok := i.values[t]
	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	tag := strconv.Itoa(index)
	m[tag] = reflect.ValueOf(value)

	return i
}

func (i *injector) SetIndex(t reflect.Type, v reflect.Value, index int) TypeIndexMapper {
	m, ok := i.values[t]
	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	tag := strconv.Itoa(index)
	m[tag] = v

	return i
}

func (i *injector) GetIndex(t reflect.Type, index int) reflect.Value {
	v := reflect.ValueOf(nil)

	m, ok := i.values[t]
	if ok {
		tag := strconv.Itoa(index)
		v := m[tag]

		if v.IsValid() {
			return v
		}

		if t.Kind() == reflect.Interface {
			vt := v.Type()
			if vt.Implements(t) {
				return v
			}
		}
	}

	if !v.IsValid() && i.parent != nil {
		v = i.parent.GetIndex(t, index)
	}

	return v
}

func (i *injector) MapTag(value interface{}, tag string) TypeIndexMapper {
	t := reflect.TypeOf(value)

	m, ok := i.values[t]

	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	m[tag] = reflect.ValueOf(value)

	return i
}

func (i *injector) MapTagTo(value interface{}, interfacePtr interface{}, tag string) TypeIndexMapper {
	t := InterfaceOfPtr(interfacePtr)

	m, ok := i.values[t]
	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	m[tag] = reflect.ValueOf(value)

	return i
}

func (i *injector) SetTag(t reflect.Type, v reflect.Value, tag string) TypeIndexMapper {
	m, ok := i.values[t]
	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	m[tag] = v

	return i
}

func (i *injector) GetTag(t reflect.Type, tag string) reflect.Value {
	v := reflect.ValueOf(nil)

	m, ok := i.values[t]
	if ok {
		v := m[tag]

		if v.IsValid() {
			return v
		}

		if t.Kind() == reflect.Interface {
			vt := v.Type()
			if vt.Implements(t) {
				return v
			}
		}
	}

	if !v.IsValid() && i.parent != nil {
		v = i.parent.GetTag(t, tag)
	}

	return v
}
