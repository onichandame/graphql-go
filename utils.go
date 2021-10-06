package graphqlgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
)

const TAG_KEY = "graphqlgo"

func getTags(field *reflect.StructField) (res map[string]string) {
	res = make(map[string]string)
	if tags, err := structtag.Parse(string(field.Tag)); err != nil {
		panic(err)
	} else {
		if tag, err := tags.Get(TAG_KEY); err == nil {
			rawTags := make([]string, 0)
			rawTags = append(rawTags, tag.Name)
			rawTags = append(rawTags, tag.Options...)
			for _, rawTag := range rawTags {
				tuple := strings.Split(rawTag, "=")
				if len(tuple) > 2 {
					panic(fmt.Errorf("tag value cannot include '=': %s", rawTag))
				}
				key := strings.TrimSpace(tuple[0])
				var value string
				if len(tuple) > 1 {
					value = tuple[1]
				}
				res[key] = value
			}
		}
	}
	return res
}

func unwrapPtr(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return unwrapPtr(t.Elem())
	} else {
		return t
	}
}
