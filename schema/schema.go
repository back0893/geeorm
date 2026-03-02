package schema

import (
	"geemod/dialect"
	"reflect"
)

type Field struct {
	Name string
	Tag  string
	Type string
}
type Schema struct {
	Model      any
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldsMap  map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldsMap[name]
}

func Parse(model any, d dialect.Dialect) *Schema {
	//获得结构体
	rv := reflect.Indirect(reflect.ValueOf(model))
	if rv.Kind() != reflect.Struct {
		panic("model must be struct")
	}
	tv := rv.Type()
	schema := Schema{
		Model:     model,
		Name:      tv.Name(),
		FieldsMap: make(map[string]*Field),
	}
	for i := 0; i < tv.NumField(); i++ {
		p := tv.Field(i)
		if !p.Anonymous && p.IsExported() {
			filed := Field{
				Name: p.Name,
				Tag:  p.Tag.Get("geeorm"),
				Type: d.DataTypeOf(reflect.Indirect(rv.Field(i))),
			}
			schema.Fields = append(schema.Fields, &filed)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.FieldsMap[p.Name] = &filed
		}

	}
	return &schema
}

func (s *Schema) RecordValues(dest any) []any {
	refDest := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []any
	for i := 0; i < len(s.Fields); i++ {
		fieldValues = append(fieldValues, refDest.Field(i).Interface())
	}
	return fieldValues
}
