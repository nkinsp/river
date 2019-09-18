package river

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type Convert struct {}

func convertName(name string) string  {
	str := ""
	for i,s:= range name{
		if i == 0 {
			str+=strings.ToLower(string(s))
		}else {
			str+=string(s)
		}
	}
	return str

}

func stringConvertTo(v string,t reflect.Type) (reflect.Value,error)  {
	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(v),nil
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		i,err := strconv.ParseInt(v,10,64)
		if err != nil || reflect.Zero(t).OverflowInt(i) {
			return reflect.ValueOf(nil),err
		}
		return reflect.ValueOf(i).Convert(t),nil
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		i,err := strconv.ParseUint(v,10,64)
		if err != nil || reflect.Zero(t).OverflowUint(i){
			return reflect.ValueOf(nil),err
		}
		return reflect.ValueOf(i).Convert(t),nil
	case reflect.Float32,reflect.Float64:
		i,err := strconv.ParseFloat(v,64)
		if err != nil || reflect.Zero(t).OverflowFloat(i){
			return reflect.ValueOf(nil),err
		}
		return reflect.ValueOf(i).Convert(t),nil
	case reflect.Bool:
		b,err := strconv.ParseBool(v)
		if err != nil {
			return reflect.ValueOf(nil),err
		}
		return reflect.ValueOf(b),nil
	case reflect.Interface:
		fmt.Println(t.Name(),"===")
		return reflect.ValueOf(v).Convert(t),nil

	}


	return reflect.ValueOf(nil),errors.New("No match "+t.Kind().String())
}

func  stringTo(v []string,t reflect.Value) error {

	switch t.Kind() {
	case reflect.Slice:
		var values []reflect.Value
		for _,s:= range v{
			value,err :=stringConvertTo(s,t.Type().Elem())
			if err != nil {
				return  err
			}
			values = append(values,value)
		}
		t.Set(reflect.Append(t,values...))

	default:
		value,err :=stringConvertTo(v[0],t.Type())
		if err != nil {
			return  err
		}
		t.Set(value)
	}


	return nil

}

func convertFormTo(form url.Values, v interface{}) error  {

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New(rv.String()+" must be a pointer")
	}
	reflectValue := rv.Elem()
	reflectType := reflect.TypeOf(v).Elem()
	num := reflectType.NumField()
	for i:=0;i < num;i++ {
		fieldValue :=reflectValue.Field(i)
		fieldType :=reflectType.Field(i)
		key := fieldType.Tag.Get("form")
		if key == "" {
			key = convertName(fieldType.Name)
		}
		value,exists := form[key]
		if !exists {
			continue
		}
		err := stringTo(value,fieldValue)
		if err != nil {
			log.Println("[River] ConvertForm field ",key," error ",err.Error())
		}
	}
	return nil
}

func ConvertToMap(v interface{}) (map[string]interface{})  {

	data := map[string]interface{}{}
	value := reflect.ValueOf(v)
	types := reflect.TypeOf(v)
	num := value.NumField()
	for i :=0; i<num;i++ {
		fieldValue :=value.Field(i)
		fieldType :=	types.Field(i)
		data[convertName(fieldType.Name)] = fieldValue.Interface()

	}

	return data

}

func mapping(fieldValue reflect.Value,value reflect.Value) error {


	fieldValue.Set(value)

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value.String())
	case reflect.Slice:
	case reflect.Map:
		//fieldValue.set
	case reflect.Struct:
	case reflect.Bool:
		fieldValue.SetBool(value.Bool())
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		//str
		//fieldValue.SetInt(value.)
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
	}

	return  nil
}

func ConvertMapTo(data map[string]interface{},v interface{}) error  {

	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.IsNil() {
		return errors.New(reflectValue.String()+" must be a pointer")
	}
	reflectValue =reflectValue.Elem()
	reflectType := reflect.TypeOf(v).Elem()
	num := reflectValue.NumField()
	for i := 0;i < num; i++{
		fieldValue := reflectValue.Field(i)
		fieldType :=reflectType.Field(i)
		key := fieldType.Tag.Get("field")
		fmt.Println("fieldType::",fieldValue.Kind(),fieldValue.Type())
		if fieldValue.Kind() == reflect.Slice {
			//rn :=reflect.New(fieldValue.Type())
			fmt.Println("rn==>",fieldValue.Type().Kind())
		}

		if key == "" {
			key = convertName(fieldType.Name)
		}
		value,exists := data[key]
		if !exists {
			continue
		}
		mapping(fieldValue,reflect.ValueOf(value))
	}

	return nil
}