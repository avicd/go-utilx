package refx

import (
	"fmt"
	"github.com/avicd/go-utilx/conv"
	"github.com/avicd/go-utilx/logx"
	"math/big"
	"reflect"
)

const (
	CmpNeq = -2
	CmpLss = -1
	CmpEq  = 0
	CmpGtr = 1
)

func Zero() reflect.Value {
	return reflect.Value{}
}

func ZeroOf(vl any) reflect.Value {
	return reflect.New(TypeOf(vl)).Elem()
}

func ValueOf(vl any) reflect.Value {
	switch tmp := vl.(type) {
	case reflect.Value:
		return tmp
	case *reflect.Value:
		return *tmp
	default:
		return reflect.ValueOf(vl)
	}
}

func TypeOf(vl any) reflect.Type {
	var value reflect.Value
	switch tmp := vl.(type) {
	case reflect.Type:
		return tmp
	case reflect.Value:
		value = tmp
	case *reflect.Value:
		value = *tmp
	}
	if value.IsValid() {
		return value.Type()
	} else {
		return reflect.TypeOf(vl)
	}
}

func KindOf(vl any) reflect.Kind {
	switch tmp := vl.(type) {
	case reflect.Value:
		return tmp.Kind()
	case *reflect.Value:
		return tmp.Kind()
	case reflect.Type:
		if tmp != nil {
			return tmp.Kind()
		} else {
			return reflect.Invalid
		}
	default:
		if tp := reflect.TypeOf(vl); tp != nil {
			return tp.Kind()
		}
	}
	return reflect.Invalid
}

func IndirectKind(val any) reflect.Kind {
	if IsPointer(val) || IsInterface(val) {
		switch tmp := val.(type) {
		case reflect.Value:
			return KindOf(IndirectType(tmp))
		case *reflect.Value:
			return KindOf(IndirectType(*tmp))
		case reflect.Type:
			return KindOf(IndirectType(tmp))
		default:
			return KindOf(IndirectType(val))
		}
	}
	return KindOf(val)
}

func Indirect(vl any) reflect.Value {
	buf := ValueOf(vl)
	for buf.Kind() == reflect.Pointer || buf.Kind() == reflect.Interface {
		buf = buf.Elem()
	}
	return buf
}

func IndirectType(vl any) reflect.Type {
	buf := TypeOf(vl)
	for buf != nil {
		if buf.Kind() == reflect.Pointer {
			buf = buf.Elem()
		} else {
			break
		}
	}
	return buf
}

func NewOf(vl any) reflect.Value {
	target := TypeOf(vl)
	if IsInvalid(target) {
		return Zero()
	}
	if target.Kind() == reflect.Pointer {
		p := reflect.New(target)
		p.Elem().Set(NewOf(target.Elem()).Addr())
		return p.Elem()
	} else if target.Kind() == reflect.Map {
		return reflect.MakeMap(target)
	} else {
		return reflect.New(target).Elem()
	}
}

func AddrAbleOf(vl any) (reflect.Value, bool) {
	buf := ValueOf(vl)
	for (IsPointer(buf) || IsInterface(buf)) && !buf.CanAddr() {
		buf = buf.Elem()
	}
	if buf.CanAddr() {
		return buf, true
	}
	return reflect.Value{}, false
}

func Assign(dest any, vl any) bool {
	target, ok := AddrAbleOf(dest)
	if ok {
		target.Set(NewOf(target))
	} else {
		logx.Fatal("dest is unaddressable")
	}
	if IsNil(vl) {
		return false
	}
	value := ValueOf(vl)
	if IsPointer(target) {
		if target.Type() == value.Type() {
			target.Set(value)
			return true
		} else if target.Type().Elem() == value.Type() {
			buf := NewOf(target)
			buf.Elem().Set(value)
			target.Set(buf)
			return true
		}
	} else if IsInterface(target) {
		if IsPointer(value) {
			if value.CanConvert(TypeOf(target.Type())) {
				target.Set(value)
				return true
			}
		} else {
			pointer := reflect.New(TypeOf(value))
			if pointer.CanConvert(target.Type()) {
				pointer.Elem().Set(value)
				target.Set(pointer)
				return true
			}
		}
	}
	el := Indirect(value)
	if target.Type() == el.Type() {
		target.Set(el)
		return true
	}
	failed := false
	if IsNumber(target) {
		if IsString(el) {
			if IsGeneralInt(target) {
				num := conv.ParseInt(el.String())
				if target.CanInt() {
					target.SetInt(num)
				} else {
					target.SetUint(uint64(num))
				}
			} else {
				target.SetFloat(conv.ParseFloat(el.String()))
			}
			return true
		}

		if target.CanInt() {
			if el.CanInt() {
				target.SetInt(el.Int())
			} else if el.CanUint() {
				target.SetInt(int64(el.Uint()))
			} else if el.CanFloat() {
				target.SetInt(int64(el.Float()))
			} else {
				failed = true
			}
		} else if target.CanUint() {
			if el.CanUint() {
				target.SetUint(el.Uint())
			} else if el.CanInt() {
				target.SetUint(uint64(el.Int()))
			} else if el.CanFloat() {
				target.SetUint(uint64(el.Float()))
			} else {
				failed = true
			}
		} else if target.CanFloat() {
			if el.CanFloat() {
				target.SetFloat(el.Float())
			} else if el.CanInt() {
				target.SetFloat(float64(el.Int()))
			} else if el.CanUint() {
				target.SetFloat(float64(el.Uint()))
			} else {
				failed = true
			}
		}
	} else if IsString(target) {
		if IsString(el) {
			target.SetString(el.String())
		} else {
			target.SetString(fmt.Sprintf("%v", vl))
		}
	} else if IsBool(target) {
		if IsBool(el) {
			target.SetBool(el.Bool())
		} else {
			target.SetBool(!el.IsZero())
		}
	} else {
		failed = true
	}
	if failed {
		logx.Fatalf("can't assign %s to %s", TypeOf(vl), TypeOf(dest))
	}
	return true
}

func Clone(dest any, src any) {
	var target reflect.Value
	if buf, ok := AddrAbleOf(dest); ok {
		buf.Set(NewOf(buf))
		target = Indirect(buf)
	} else {
		logx.Fatal("dest is unaddressable")
	}
	el := Indirect(src)
	if IsNil(el) {
		return
	}
	if el.Type() != target.Type() {
		logx.Fatalf("can't clone %s into %s", el.Type(), target.Type())
	}
	if IsBasic(target) {
		target.Set(el)
		return
	}
	switch el.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < el.Len(); i++ {
			sf := el.Index(i)
			if !IsBasic(sf) {
				newVal := NewOf(el.Type().Elem())
				Clone(newVal.Addr(), sf)
				sf = newVal
			}
			newTarget := reflect.Append(target, sf)
			target.Set(newTarget)
		}
	case reflect.Struct:
		for i := 0; i < el.Type().NumField(); i++ {
			tf := el.Type().Field(i)
			if !tf.IsExported() {
				continue
			}
			sf := el.Field(i)
			if !IsBasic(tf.Type) {
				target.Field(i).Set(sf)
			} else {
				newVal := NewOf(tf.Type)
				Clone(newVal.Addr(), sf)
				sf = newVal
			}
			target.Field(i).Set(sf)
		}
	case reflect.Map:
		for itr := el.MapRange(); itr.Next(); {
			sf := itr.Value()
			if IsBasic(sf) {
				target.SetMapIndex(itr.Key(), sf)
			} else {
				var newVal reflect.Value
				if IsInterface(el.Type().Elem()) {
					newVal = NewOf(sf.Elem())
				} else {
					newVal = NewOf(el.Type().Elem())
				}
				Clone(newVal.Addr(), sf)
				target.SetMapIndex(itr.Key(), newVal)
			}
		}
	}
}

func Merge(dest any, src any) {
	if IsNil(src) {
		return
	}
	to := Indirect(dest)
	if IsBasic(to) || IsList(to) || IsNil(to) {
		if target, ok := AddrAbleOf(dest); ok {
			if IsNil(target) {
				target.Set(NewOf(target))
			}
			to = Indirect(target)
		} else {
			logx.Fatal("dest is unaddressable")
		}
	}
	from := Indirect(src)

	if !IsList(to) && !IsList(from) && to.Type() != from.Type() {
		logx.Fatalf("can't merge %s into %s", TypeOf(src), TypeOf(dest))
	}
	switch to.Kind() {
	case reflect.Slice:
		buf := NewOf(to)
		buf.Set(to)
		for i := 0; i < from.Len(); i++ {
			sf := from.Index(i)
			buf = reflect.Append(buf, sf)
		}
		to.Set(buf)
	case reflect.Array:
		for i := 0; i < from.Len() && i < to.Len(); i++ {
			sf := from.Index(i)
			if IsInterface(sf) {
				sf = sf.Elem()
			}
			if !sf.IsZero() {
				to.Index(i).Set(sf)
			}
		}
	case reflect.Struct:
		for i := 0; i < from.Type().NumField(); i++ {
			tf := from.Type().Field(i)
			if !tf.IsExported() {
				continue
			}
			left := to.Field(i)
			right := from.Field(i)
			if IsBasic(tf.Type) {
				if right.IsValid() && !right.IsZero() {
					left.Set(right)
				}
			} else if !right.IsZero() {
				if IsNil(left) {
					left.Set(NewOf(tf.Type))
				}
				Merge(left.Addr(), right)
			}
		}
	case reflect.Map:
		for itr := from.MapRange(); itr.Next(); {
			left := to.MapIndex(itr.Key())
			right := itr.Value()
			if IsInterface(left) {
				left = left.Elem()
			}
			if IsInterface(right) {
				right = right.Elem()
			}
			if left.IsValid() {
				if !IsBasic(left) && left.Type() == right.Type() {
					var newVal reflect.Value
					if IsInterface(from.Type().Elem()) {
						newVal = NewOf(left)
					} else {
						newVal = NewOf(from.Type().Elem())
					}
					newVal.Set(left)
					Merge(newVal.Addr(), right)
					to.SetMapIndex(itr.Key(), newVal)
				} else if !right.IsZero() {
					to.SetMapIndex(itr.Key(), right)
				}
			} else {
				to.SetMapIndex(itr.Key(), right)
			}
		}
	default:
		if !from.IsZero() {
			to.Set(from)
		}
	}
}

func FieldOf(src any, props ...any) (reflect.Value, bool) {
	vl := ValueOf(src)
	if IsBasic(vl) || len(props) < 1 || IsNil(vl) {
		return Zero(), false
	}
	buf := vl
	for _, key := range props {
		el := Indirect(buf)
		switch el.Kind() {
		case reflect.Map:
			index := AsOf(el.Type().Key(), key)
			buf = el.MapIndex(index)
			if IsInterface(buf) {
				buf = buf.Elem()
			}
		case reflect.Struct:
			name := AsString(key)
			if sf, ok := el.Type().FieldByName(name); ok && sf.IsExported() {
				buf = el.FieldByName(name)
			} else {
				return Zero(), false
			}
		case reflect.Slice, reflect.Array:
			index := AsInt(key)
			if index >= 0 && index < el.Len() {
				buf = el.Index(index)
			} else {
				return Zero(), false
			}
		default:
			return Zero(), false
		}
		if !buf.IsValid() || !buf.CanInterface() {
			return Zero(), false
		}
	}
	return buf, true
}

func TypeOfField(vl any, props ...any) (reflect.Type, bool) {
	if _, ok := vl.(reflect.Type); !ok && !IsNil(vl) {
		if sf, exist := FieldOf(vl, props...); exist {
			return TypeOf(sf), true
		}
	}
	buf := TypeOf(vl)
	if len(props) < 1 || buf == nil {
		return nil, false
	}
	for _, key := range props {
		buf = IndirectType(buf)
		switch buf.Kind() {
		case reflect.Struct:
			name := AsString(key)
			if sf, ok := buf.FieldByName(name); ok {
				buf = sf.Type
			} else {
				return nil, false
			}
		case reflect.Array, reflect.Slice, reflect.Map:
			buf = buf.Elem()
		default:
			return nil, false
		}
	}
	return buf, true
}

func TypeOfId(vl any, ident string) (reflect.Type, bool) {
	props := AsList(conv.StrToArr(ident, "."))
	return TypeOfField(vl, props...)
}

func PropOf(src any, props ...any) (any, bool) {
	vl := ValueOf(src)
	if IsBasic(vl) || len(props) < 1 || IsNil(vl) {
		return nil, false
	}
	if value, ok := FieldOf(vl, props...); ok {
		return value.Interface(), true
	} else {
		return nil, false
	}
}

func PropOfId(src any, ident string) (any, bool) {
	if IsBasic(src) || IsNil(src) {
		return nil, false
	}
	props := AsList(conv.StrToArr(ident, "."))
	if val, ok := PropOf(src, props...); ok {
		return val, true
	} else if IsMap(src) && len(props) > 1 {
		return PropOf(src, ident)
	} else {
		return nil, false
	}
}

func Set(dest any, vl any, props ...any) bool {
	if IsBasic(dest) {
		return Assign(dest, vl)
	}
	value := ValueOf(vl)
	buf := Indirect(dest)
	if IsNil(buf) || IsArray(dest) {
		if target, ok := AddrAbleOf(dest); !ok {
			logx.Fatal("dest is unaddressable")
		} else {
			buf = target
		}
	}
	for i, key := range props {
		if IsNil(buf) {
			buf.Set(NewOf(buf))
		}
		el := Indirect(buf)
		switch el.Kind() {
		case reflect.Map:
			index := AsOf(el.Type().Key(), key)
			buf = el.MapIndex(index)
			if len(props)-i > 1 {
				if !buf.IsValid() {
					if IsInterface(el.Type().Elem()) {
						next := props[i+1]
						if IsNumber(next) || IsString(next) {
							buf = NewOf(TMapStrAny)
						} else {
							buf = NewOf(reflect.MapOf(TypeOf(next), el.Type().Elem()))
						}
					} else {
						buf = NewOf(el.Type().Elem())
					}
					el.SetMapIndex(index, buf)
					continue
				}
			} else {
				el.SetMapIndex(index, value)
				break
			}
		case reflect.Struct:
			name := AsString(key)
			if tf, ok := el.Type().FieldByName(name); ok && tf.IsExported() {
				buf = el.FieldByName(name)
				if IsPointer(buf) && buf.IsNil() {
					buf.Set(NewOf(tf.Type))
				}
			} else {
				return false
			}
		case reflect.Slice, reflect.Array:
			index := AsInt(key)
			if index >= 0 && index < el.Len() {
				buf = el.Index(index)
			} else {
				return false
			}
		}
		if i == len(props)-1 {
			if buf.IsValid() {
				Assign(buf.Addr(), value)
			} else {
				return false
			}
		}
	}
	return true
}

func SetById(dest any, vl any, ident string) bool {
	props := AsList(conv.StrToArr(ident, "."))
	if ok := Set(dest, vl, props...); ok {
		return true
	} else if IsMap(dest) && len(props) > 1 {
		return Set(dest, vl, ident)
	} else {
		return false
	}
}

func MethodOf(target any, props ...any) (reflect.Value, bool) {
	vl := ValueOf(target)
	if !vl.IsValid() || IsBasic(vl) || len(props) < 1 {
		return Zero(), false
	}
	owner := vl
	level := len(props)
	if level > 1 {
		if sf, ok := FieldOf(owner, props[:level-1]...); ok {
			owner = sf
		} else {
			return Zero(), false
		}
	}
	name := AsString(props[level-1])
	method := owner.MethodByName(name)
	if method.IsValid() && method.CanInterface() {
		return method, true
	}
	return Zero(), false
}

func MethodOfId(target any, ident string) (reflect.Value, bool) {
	props := AsList(conv.StrToArr(ident, "."))
	return MethodOf(target, props...)
}

func ForEach(target any, fn func(key any, val any)) {
	if IsBasic(target) || IsNil(target) {
		return
	}
	vl := ValueOf(target)
	el := Indirect(vl)
	switch el.Kind() {
	case reflect.Map:
		for itr := el.MapRange(); itr.Next(); {
			fn(itr.Key().Interface(), itr.Value().Interface())
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < el.Len(); i++ {
			fn(i, el.Index(i).Interface())
		}
	case reflect.Struct:
		for i := 0; i < el.Type().NumField(); i++ {
			tf := el.Type().Field(i)
			field := el.Field(i)
			if !tf.IsExported() {
				continue
			}
			fn(tf.Name, field.Interface())
		}
	}
}

func Cmp(x any, y any) int {
	if IsNumber(x) && IsNumber(y) {
		return big.NewFloat(AsFloat64(x)).Cmp(big.NewFloat(AsFloat64(y)))
	} else if IsString(x) && IsString(y) {
		str1 := AsString(x)
		str2 := AsString(y)
		if str1 == str2 {
			return CmpEq
		} else if str1 > str2 {
			return CmpGtr
		} else {
			return CmpLss
		}
	} else if reflect.DeepEqual(x, y) {
		return CmpEq
	}
	return CmpNeq
}
