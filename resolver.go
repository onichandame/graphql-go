package graphql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

const (
	QUERY_PREFIX    = "Query"
	MUTATION_PREFIX = "Mutation"
)

type ResolverFunc func(ctx *gin.Context, args interface{}) interface{}

type Resolvers struct {
	Query    graphql.Fields
	Mutation graphql.Fields
}

func GetResolverFromStruct(instance interface{}, objType *graphql.Object) (resolvers *Resolvers) {
	resolvers = new(Resolvers)
	resolvers.Query = make(graphql.Fields)
	resolvers.Mutation = make(graphql.Fields)
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("resolver must be defined by a struct but %s is not", t.Name()))
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if strings.HasPrefix(m.Name, QUERY_PREFIX) {
			name := strings.Replace(m.Name, QUERY_PREFIX, "", 1)
			args, fn := getResolverFunc(&m)
			resolvers.Query[name] = &graphql.Field{
				Type:    objType,
				Args:    args,
				Resolve: fn,
			}
		} else if strings.HasPrefix(m.Name, MUTATION_PREFIX) {
			name := strings.Replace(m.Name, MUTATION_PREFIX, "", 1)
			args, fn := getResolverFunc(&m)
			resolvers.Mutation[name] = &graphql.Field{
				Type:    objType,
				Args:    args,
				Resolve: fn,
			}
		}
	}
	return resolvers
}

func validateFunc(fn *reflect.Method) {
	var standard ResolverFunc
	if reflect.TypeOf(standard) != fn.Type {
		panic(fmt.Errorf("resolver function must be of ResolverFunc type, %s is not", reflect.TypeOf(fn).Name()))
	}
}

func getResolverFunc(m *reflect.Method) (args graphql.FieldConfigArgument, fn graphql.FieldResolveFn) {
	args = make(graphql.FieldConfigArgument)
	validateFunc(m)
	ctxType := m.Type.In(0)
	argsType := m.Type.In(1)
	fn = func(p graphql.ResolveParams) (res interface{}, err error) {
		return res, err
	}
	return args, fn
}
