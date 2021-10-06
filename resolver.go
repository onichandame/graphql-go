package graphqlgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	goutils "github.com/onichandame/go-utils"
)

const (
	QUERY_PREFIX        = "Query"
	MUTATION_PREFIX     = "Mutation"
	SUBSCRIPTION_PREFIX = "Subscription" // TODO
)

type ResolverFunc func(input interface{}, params graphql.ResolveParams) interface{}

type Resolver struct {
	Query    graphql.Fields
	Mutation graphql.Fields
}

func GetResolverFromStruct(instance interface{}, inputType *graphql.Object) (resolver *Resolver) {
	resolver = new(Resolver)
	resolver.Query = make(graphql.Fields)
	resolver.Mutation = make(graphql.Fields)
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("resolver must be defined by a struct but %s is not", t.Name()))
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		var fields *graphql.Fields
		var name string
		if strings.HasPrefix(m.Name, QUERY_PREFIX) {
			name = strings.Replace(m.Name, QUERY_PREFIX, "", 1)
			fields = &resolver.Query
		} else if strings.HasPrefix(m.Name, MUTATION_PREFIX) {
			name = strings.Replace(m.Name, MUTATION_PREFIX, "", 1)
			fields = &resolver.Mutation
		} else {
			continue
		}
		validateFunc(&m)
		argsObject := GetObjectFromStruct(reflect.New(m.Type.In(0)).Interface(), INPUT)
		args := make(graphql.FieldConfigArgument)
		for _, field := range argsObject.Fields {
			args[field.Name] = &graphql.ArgumentConfig{
				Type: field.Type,
			}
		}
		(*fields)[name] = &graphql.Field{
			Type: inputType,
			Args: args,
			Resolve: func(p graphql.ResolveParams) (res interface{}, err error) {
				defer goutils.RecoverToErr(&err)
				params := make([]reflect.Value, 0)
				argsType := m.Type.In(0)
				if argsType.Kind() == reflect.Ptr {
					argsType = argsType.Elem()
				}
				args := reflect.New(argsType).Interface()
				if err := mapstructure.Decode(p.Args, args); err != nil {
					panic(err)
				}
				params = append(params, reflect.ValueOf(args))
				params = append(params, reflect.ValueOf(p))
				res = m.Func.Call(params)[0].Interface()
				return res, err
			},
		}
	}
	return resolver
}

func validateFunc(fn *reflect.Method) {
	var standard ResolverFunc
	if reflect.TypeOf(standard) != fn.Type {
		panic(fmt.Errorf("resolver function must be of ResolverFunc type, %s is not", reflect.TypeOf(fn).Name()))
	}
}

func getResolverFunc(m *reflect.Method) (args graphql.FieldConfigArgument, fn graphql.FieldResolveFn) {
	validateFunc(m)
	inputType := m.Type.In(0)
	for i := 0; i < inputType.NumField(); i++ {
		fieldType := inputType.Field(i)
		tags := getTags(&fieldType)
		name := fieldType.Name
		var argType graphql.Input
		if customName, ok := tags["name"]; ok {
			name = customName
		}
		args[name] = &graphql.ArgumentConfig{
			Type: argType,
		}

	}
	args = make(graphql.FieldConfigArgument)
	fn = func(p graphql.ResolveParams) (res interface{}, err error) {
		runHandler := func() interface{} {
			defer goutils.RecoverToErr(&err)
			args := make([]reflect.Value, 0)
			args[0] = reflect.ValueOf(p)
			return m.Func.Call(args)[0].Interface()
		}
		res = runHandler()
		return res, err
	}
	return args, fn
}
