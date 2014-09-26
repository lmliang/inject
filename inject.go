package inject

import (
	"fmt"
	"reflect"
	"strconv"
)

// 函数调用
type Invoke interface {
	// 调用函数，以reflect.Value方式返回参数及错误信息，如果参数不是函数则panic
	Invoke(interface{}) ([]reflect.Value, error)
}

// 赋值
type Assign interface {
	// 结构体字段赋值
	AssignField(interface{}) error
}

// 按照index进行Type映射
type TypeIndexMapper interface {
	MapIndex(int, interface{}) TypeIndexMapper

	MapIndexTo(int, interface{}, interface{}) TypeIndexMapper

	SetIndex(int, reflect.Type, reflect.Value) TypeIndexMapper

	GetIndex(int, reflect.Type) reflect.Value
}

// 按照tag进行Type映射
type TypeTagMapper interface {
	MapTag(string, interface{}) TypeTagMapper

	MapTagTo(string, interface{}, interface{}) TypeTagMapper

	SetTag(string, reflect.Type, reflect.Value) TypeTagMapper

	GetTag(string, reflect.Type) reflect.Value
}

type TypeMapper interface {
	TypeIndexMapper
	TypeTagMapper
}

type Injector interface {
	Invoke
	Assign
	TypeMapper

	SetParent(Injector)
}

// 获得指向interface指针的类型
func InterfaceOfPtr(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic("Inject.Interface: Expect a pointer to an interface.")
	}

	return t
}

type injector struct {
	// 按照reflect.Type类型为键，值按照标签存储
	// 函数多参数类型相同时，标签为参数序号；
	// 结构体字段赋值时标签为字段名称（必须为可导出字段）
	values map[reflect.Type]map[string]reflect.Value
	parent Injector
}

func (i *injector) Invoke(fn interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(fn)

	in := make([]reflect.Value, t.NumIn())

	for n := 0; n < t.NumIn(); n++ {
		tp := t.In(n)
		arg := i.GetIndex(n, tp)
		if !arg.IsValid() {
			return nil, fmt.Errorf("Inject-Invoke: call func %s not found param(%d) type %#v", t.Name(), n, tp.Name())
		}
		in[n] = arg
	}

	return reflect.ValueOf(fn).Call(in), nil
}

func (i *injector) AssignField(st interface{}) error {
	v := reflect.ValueOf(st)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("Inject-AssignField: expect a struct param.")
	}

	t := v.Type()

	for n := 0; n < t.NumField(); n++ {
		f := v.Field(n)
		sf := t.Field(n)

		val := i.GetTag(sf.Name, f.Type())
		if val.IsValid() {
			f.Set(val)
		}
	}

	return nil
}

func (i *injector) SetParent(parent Injector) {
	i.parent = parent
}

func (i *injector) MapIndex(index int, value interface{}) TypeIndexMapper {
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

func (i *injector) MapIndexTo(index int, value interface{}, interfacePtr interface{}) TypeIndexMapper {
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

func (i *injector) SetIndex(index int, t reflect.Type, v reflect.Value) TypeIndexMapper {
	m, ok := i.values[t]
	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	tag := strconv.Itoa(index)
	m[tag] = v

	return i
}

func (i *injector) GetIndex(index int, t reflect.Type) reflect.Value {
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
		v = i.parent.GetIndex(index, t)
	}

	return v
}

func (i *injector) MapTag(tag string, value interface{}) TypeTagMapper {
	t := reflect.TypeOf(value)

	m, ok := i.values[t]

	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	m[tag] = reflect.ValueOf(value)

	return i
}

func (i *injector) MapTagTo(tag string, value interface{}, interfacePtr interface{}) TypeTagMapper {
	t := InterfaceOfPtr(interfacePtr)

	m, ok := i.values[t]
	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	m[tag] = reflect.ValueOf(value)

	return i
}

func (i *injector) SetTag(tag string, t reflect.Type, v reflect.Value) TypeTagMapper {
	m, ok := i.values[t]
	if !ok {
		m = make(map[string]reflect.Value)
		i.values[t] = m
	}

	m[tag] = v

	return i
}

func (i *injector) GetTag(tag string, t reflect.Type) reflect.Value {
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
		v = i.parent.GetTag(tag, t)
	}

	return v
}

func New() Injector {
	return &injector{values: make(map[reflect.Type]map[string]reflect.Value)}
}
