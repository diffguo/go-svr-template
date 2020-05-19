package tools

import (
	"fmt"
	"reflect"
)

///////////////////////////////////////////////////////////////////////////////////////
/*
// 为了高效、便捷的为Slice中的对象结构中的成员进行赋值操作（待赋值的数据可能在数据库中），从而实现该方法
//
// 如有结构：
type T struct {
	ID int64
	Name string
	Member MObj
}

生成了对应的Slice：objs := []T{T{...},T{...},T{...},....}
需要填充其中成员的Member域的值，使用如下结构的对象：
type MObj struct {
	Id int
	Name string
}

如果MObj的数据在数据库中，通常的做法是：
for _, t := range objs {
	1. fetch mobj from db by t.Id
	2. t.Member = mobj
}

这样的缺点是数据库访问次数多，代码量大，不易维护。

使用DataBox后的方式是：

err := NewDataBox(dataSlice).KeyField("ID").JoinByMap(func(keywords interface{}) interface{} {
		members := make(map[int64]Obj) // Obj不能是指针
		ids := keywords.([]int64)
		// 一次性从数据库中得到所有的Mobjs
		for _, id := range ids {
			members[id] = MObj{....}
		}

		return members
	}).SaveToField("Member")
//
*/
///////////////////////////////////////////////////////////////////////////////////////

type DataBox struct {
	// 主数据, 为slice
	data interface{}

	// 数据基本类型
	dataValue reflect.Value

	// ptr type of the element in the data
	elemPtrType reflect.Type

	// 主数据Map，主key(elemPtrType)指向主数据的成员。key为主数据Slice中的Element中的keyFieldName对应的值，Value为主数据Slice中的Element的引用
	elemMap reflect.Value

	// 关联Map，主key(elemPtrType)指向返回的数据对象
	relatedMap reflect.Value

	// 内部错误传递，用户不用关心该成员
	err error
}

// 用户需要实现的函数：通过key列表来获取对象
// key列表时data里面通过keyfield抽取的key列表
// 第一个返回的map的key为传入的key；map里面的value可以是查到的单个对象，也可以是一个可以对应的对象slice
// 与FetchObjsByKeysRetSliceFunc相比，性能更高，且子结构可以不带父对象的ID，缺点是需要多写点代码
type FetchObjsByKeysRetMapFunc func(keyList interface{}) (key2ObjsMap interface{}) // map[key][]obj  or  map[key]obj  or map[key][]*obj  or  map[key]*obj

// 第二个返回对象的切片
// 与FetchObjsByKeysRetMapFunc相比，代码更简洁，但是子结构必须带父对象的ID，以便系统进行关联，缺点是：由于有内存拷贝，性能相对较差
type FetchObjsByKeysRetSliceFunc func(keyList interface{}) (objList interface{}) // []obj or []*obj

// 功能：生成NewDataBox对象，以便向data(必须是slice)中的成员(slice中的成员)添加指定数据，数据由FetchObjsByKeysRetMapFunc或者
// FetchObjsByKeysRetSliceFunc返回slice中的对象和map中的value可以是结构或者是结构的指针
// keyFieldName指明后续获取数据的key，后续会收集所有的key，传到FetchObjs...函数，以便集中化获取数据，提升效率
// 通过SaveToField函数指明的fieldname，得到数据填充的位置
// 其它特点：
// 1. 可以回填对象或者对象Slice
// 2. 可以多重填充
// 3. 被填充对象为指针，填充对象为结构对象 这种情况会导致异常
func NewDataBox(data interface{}) *DataBox {
	databox := &DataBox{}
	databox.data = data
	if data == nil {
		databox.err = fmt.Errorf("Empty Input Data")
		return databox
	}

	databox.dataValue = reflect.ValueOf(data)
	if !databox.dataValue.IsValid() || databox.dataValue.IsNil() {
		databox.err = fmt.Errorf("Wrong Input Data Format")
		return databox
	}

	if databox.dataValue.Kind() == reflect.Ptr {
		databox.dataValue = databox.dataValue.Elem()
		databox.data = databox.dataValue.Interface()
	}

	var elemOriginalType reflect.Type
	if databox.dataValue.Kind() == reflect.Slice || databox.dataValue.Kind() == reflect.Map { // 如果是slice类型, 取出实体类型
		elemOriginalType = databox.dataValue.Type().Elem()
	} else {
		databox.err = fmt.Errorf("Wrong Input Data's Type. Only Should Be Slice Or Map")
		return databox
	}

	// 确认真实类型
	if elemOriginalType.Kind() == reflect.Ptr {
		databox.elemPtrType = elemOriginalType
	} else {
		databox.elemPtrType = reflect.PtrTo(elemOriginalType)
	}

	if databox.elemPtrType.Elem().Kind() != reflect.Struct {
		databox.err = fmt.Errorf("Wrong Input Data's Element Type. Only Should Be Struct")
		return databox
	}

	return databox
}

func (d *DataBox) KeyField(keyFieldName string) *DataBox {
	if d.err != nil {
		return d
	}

	structField, ok := d.elemPtrType.Elem().FieldByName(keyFieldName)
	if !ok {
		d.err = fmt.Errorf("keyFieldName not in element")
		return d
	}

	d.elemMap = reflect.MakeMap(reflect.MapOf(structField.Type, reflect.SliceOf(d.elemPtrType)))

	//将原始数据, 整理进中间的比较字典
	for i := 0; i < d.dataValue.Len(); i++ {
		item := d.dataValue.Index(i)
		if item.Kind() == reflect.Ptr {
			d.storeKVIntoMap(item.Elem().FieldByIndex(structField.Index), item, &d.elemMap)
		} else {
			d.storeKVIntoMap(item.FieldByIndex(structField.Index), item.Addr(), &d.elemMap)
		}
	}

	return d
}

func (d *DataBox) JoinByMap(f FetchObjsByKeysRetMapFunc) *DataBox {
	if d.err != nil {
		return d
	}

	keyList := reflect.MakeSlice(reflect.SliceOf(d.elemMap.Type().Key()), 0, len(d.elemMap.MapKeys()))
	for _, v := range d.elemMap.MapKeys() {
		keyList = reflect.Append(keyList, v)
	}

	// 获取关联数据
	retMap := f(keyList.Interface())
	if retMap == nil {
		return nil
	}

	d.relatedMap = reflect.ValueOf(retMap)
	return d
}

func (d *DataBox) JoinByObjs(f FetchObjsByKeysRetSliceFunc, relatedfield string) *DataBox {
	if d.err != nil {
		return d
	}

	keyList := reflect.MakeSlice(reflect.SliceOf(d.elemMap.Type().Key()), 0, len(d.elemMap.MapKeys()))
	for _, v := range d.elemMap.MapKeys() {
		keyList = reflect.Append(keyList, v)
	}

	// 获取关联数据
	retObjSlice := f(keyList.Interface())
	if retObjSlice == nil {
		d.err = fmt.Errorf("User Define Function Must Returned Nil")
		return d
	}

	v := reflect.ValueOf(retObjSlice)
	t := v.Type()
	if t.Kind() != reflect.Slice {
		d.err = fmt.Errorf("User Define Function Must Return Slice")
		return d
	}

	objPtrType := t.Elem()
	if t.Elem().Kind() != reflect.Ptr {
		objPtrType = reflect.PtrTo(t.Elem())
	}

	sf, exsit := objPtrType.Elem().FieldByName(relatedfield)
	if !exsit {
		d.err = fmt.Errorf("Relatedfield: %s not In Obj Which In Return Slice", relatedfield)
		return d
	}

	if sf.Type.Kind() == reflect.Ptr {
		d.relatedMap = reflect.MakeMap(reflect.MapOf(sf.Type.Elem(), reflect.SliceOf(objPtrType)))
	} else {
		d.relatedMap = reflect.MakeMap(reflect.MapOf(sf.Type, reflect.SliceOf(objPtrType)))
	}

	// 把返回对象放入Map中
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			d.storeKVIntoMap(item.Elem().FieldByIndex(sf.Index), item, &d.relatedMap)
		} else {
			d.storeKVIntoMap(item.FieldByIndex(sf.Index), item.Addr(), &d.relatedMap)
		}
	}

	return d
}

func (d *DataBox) SaveToField(fieldName string) error {
	if d.err != nil {
		return d.err
	}

	// map -> value(slice) -> obj(pointer) -> realobj(obj)
	saveStructField, exist := d.elemMap.Type().Elem().Elem().Elem().FieldByName(fieldName)
	if !exist {
		return fmt.Errorf("%s not exsit in struct", fieldName)
	}

	// saveStructField 与 用户函数返回的map的value类型必须相同 (兼容指针情况)
	var retTypeIsPtr bool = false
	rettype := d.relatedMap.Type().Elem()
	if rettype.Kind() == reflect.Ptr {
		retTypeIsPtr = true
		rettype = d.relatedMap.Type().Elem().Elem()
	}

	var destTypeIsPtr bool = false
	dsttype := saveStructField.Type
	if dsttype.Kind() == reflect.Ptr {
		destTypeIsPtr = true
		dsttype = saveStructField.Type.Elem()
	}

	if rettype.Kind() != dsttype.Kind() {
		// 如果不相等，那么返回的可能是slice，那做特殊处理
		if rettype.Kind() == reflect.Slice {
			if rettype.Elem().Kind() == reflect.Ptr {
				if rettype.Elem().Elem().Kind() == dsttype.Kind() {
					return d.saveToFieldWithSlice(&saveStructField, destTypeIsPtr)
				}
			} else if rettype.Kind() == dsttype.Kind() {
				return d.saveToFieldWithSlice(&saveStructField, destTypeIsPtr)
			}
		}

		if rettype.Kind() == reflect.Ptr && rettype.Elem().Kind() == reflect.Slice {
			if rettype.Elem().Elem().Kind() == reflect.Ptr {
				if rettype.Elem().Elem().Elem().Kind() == dsttype.Kind() {
					return d.saveToFieldWithSlice(&saveStructField, destTypeIsPtr)
				}
			} else if rettype.Elem().Kind() == dsttype.Kind() {
				return d.saveToFieldWithSlice(&saveStructField, destTypeIsPtr)
			}
		}

		return fmt.Errorf("Type Error When SaveToField: rettype: %s != desttype: %s", rettype.String(), dsttype.String())
	}

	// 填充slice的时候，可能需要一个一个处理
	var isNeedConvertWhenElemTypeNotEquelInSlice = false
	if dsttype.Kind() == reflect.Slice && rettype.Kind() == reflect.Slice {
		if dsttype.Elem().Kind() != rettype.Elem().Kind() {
			isNeedConvertWhenElemTypeNotEquelInSlice = true
		}
	}

	// saveStructField可以是slice或者是struct或者是一个普通类型
	iter := d.relatedMap.MapRange()
	for iter.Next() {
		saveElemSlice := d.elemMap.MapIndex(iter.Key())
		if !saveElemSlice.IsValid() {
			continue
		}

		retObjValue := iter.Value()
		if isNeedConvertWhenElemTypeNotEquelInSlice {
			retObjValue = convertSlice(iter.Value(), dsttype)
		}

		// 把value放入saveElem slice中对象的fieldName处
		for i := 0; i < saveElemSlice.Len(); i++ {
			dest := saveElemSlice.Index(i).Elem().FieldByIndex(saveStructField.Index)
			if destTypeIsPtr {
				if retTypeIsPtr {
					dest.Set(retObjValue)
				} else {
					dest.Set(retObjValue.Addr())
				}
			} else {
				if retTypeIsPtr {
					dest.Set(retObjValue.Elem())
				} else {
					dest.Set(retObjValue)
				}
			}
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////
//
// 以下为内部函数
//
///////////////////////////////////////////////////////////////////////////////////////

func (d *DataBox) storeKVIntoMap(key reflect.Value, value reflect.Value, m *reflect.Value) {
	fieldDataInMapSlice := m.MapIndex(key)
	if !fieldDataInMapSlice.IsValid() {
		fieldDataInMapSlice = reflect.MakeSlice(m.Type().Elem(), 0, 1)
	}

	fieldDataInMapSlice = reflect.Append(fieldDataInMapSlice, value)
	m.SetMapIndex(key, fieldDataInMapSlice)
}

// 当使用FetchObjsByKeysRetSliceFunc函数时，返回的成员数不定，因此临时map把它们处理成slice，这里需要把slice填到Struct中去
func (d *DataBox) saveToFieldWithSlice(saveStructField *reflect.StructField, destTypeIsPtr bool) error {
	if d.err != nil {
		return d.err
	}

	var retTypeIsPtr bool = false
	rettype := d.relatedMap.Type().Elem().Elem() // map的成员为slice，再返回slice的成员
	if rettype.Kind() == reflect.Ptr {
		retTypeIsPtr = true
	}

	// saveStructField可以是slice或者是struct或者是一个普通类型
	iter := d.relatedMap.MapRange()
	for iter.Next() {
		saveElemSlice := d.elemMap.MapIndex(iter.Key())
		if !saveElemSlice.IsValid() {
			continue
		}

		if iter.Value().Len() <= 0 {
			continue
		}

		// 把value放入saveElem slice中对象的fieldName处
		for i := 0; i < saveElemSlice.Len(); i++ {
			dest := saveElemSlice.Index(i).Elem().FieldByIndex(saveStructField.Index)
			if destTypeIsPtr {
				if retTypeIsPtr {
					dest.Set(iter.Value().Index(0))
				} else {
					dest.Set(iter.Value().Index(0).Addr())
				}
			} else {
				if retTypeIsPtr {
					dest.Set(iter.Value().Index(0).Elem())
				} else {
					dest.Set(iter.Value().Index(0))
				}
			}
		}
	}

	return nil
}

// convert Slice from []Obj to []*Obj  or  from []*Obj to []Obj
func convertSlice(value reflect.Value, destSliceType reflect.Type) reflect.Value {
	if value.Type() == destSliceType {
		return value
	}

	destSlice := reflect.MakeSlice(destSliceType, 0, value.Len())

	var ptr2Struct bool = false
	if value.Type().Elem().Kind() == reflect.Ptr && destSliceType.Elem().Kind() == reflect.Struct {
		ptr2Struct = true
	}

	var struct2Prt bool = false
	if value.Type().Elem().Kind() == reflect.Struct && destSliceType.Elem().Kind() == reflect.Ptr {
		struct2Prt = true
	}

	for i := 0; i < value.Len(); i++ {
		if ptr2Struct {
			destSlice = reflect.Append(destSlice, value.Index(i).Elem())
		}

		if struct2Prt {
			destSlice = reflect.Append(destSlice, value.Index(i).Addr())
		}
	}

	return destSlice
}
