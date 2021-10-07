package graphqlgo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
	goutils "github.com/onichandame/go-utils"
)

const (
	QUERY_PREFIX        = "Query"
	MUTATION_PREFIX     = "Mutation"
	SUBSCRIPTION_PREFIX = "Subscription"
)

func getResolver(objtype reflect.Type) (query graphql.Fields, mutation graphql.Fields, subscription graphql.Fields) {
	objtype = unwrapPtr(objtype)
	query = make(graphql.Fields)
	mutation = make(graphql.Fields)
	subscription = make(graphql.Fields)
	objvalue := reflect.New(objtype)
	methods := objvalue.NumMethod()
	for i := 0; i < methods; i++ {
		operationtype := objvalue.Type().Method(i)
		var name string
		var operation string
		if strings.HasPrefix(operationtype.Name, QUERY_PREFIX) {
			operation = QUERY_PREFIX
			name = strings.Replace(operationtype.Name, QUERY_PREFIX, "", 1)
		} else if strings.HasPrefix(operationtype.Name, MUTATION_PREFIX) {
			operation = MUTATION_PREFIX
			name = strings.Replace(operationtype.Name, MUTATION_PREFIX, "", 1)
		} else if strings.HasPrefix(operationtype.Name, SUBSCRIPTION_PREFIX) {
			operation = SUBSCRIPTION_PREFIX
			name = strings.Replace(operationtype.Name, SUBSCRIPTION_PREFIX, "", 1)
		}
		if operation != "" {
			var field graphql.Field
			if operationtype.Type.NumOut() != 1 {
				panic(errors.New("the output number of an operation must be exactly 1"))
			}
			outputtype := operationtype.Type.Out(0)
			field.Type = getFieldType(outputtype, make(map[string]string), OUTPUT)
			if operationtype.Type.NumIn() < 2 || operationtype.Type.NumIn() > 3 {
				panic(fmt.Errorf("the input number of an operation must be between 1 and 2, %d is not acceptable", operationtype.Type.NumIn()))
			}
			if operationtype.Type.NumIn() > 2 {
				argstype := operationtype.Type.In(2)
				field.Args = getArgs(argstype)
			}
			field.Resolve = func(p graphql.ResolveParams) (res interface{}, err error) {
				defer goutils.RecoverToErr(&err)
				res = operationtype.Func.Call([]reflect.Value{reflect.ValueOf(nil), reflect.ValueOf(p)})[0].Interface()
				return res, err
			}
			switch operation {
			case QUERY_PREFIX:
				query[name] = &field
			case MUTATION_PREFIX:
				mutation[name] = &field
			case SUBSCRIPTION_PREFIX:
				subscription[name] = &field
			}
		}
	}
	return query, mutation, subscription
}
