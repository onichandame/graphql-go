package core

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestResolver(t *testing.T) {
	t.Run("simple query", func(t *testing.T) {
		resolver := Resolver{
			Named{Name: "Simple"},
			Describable{Description: "Simple Query"},
			[]*Query{
				{
					Operation{
						Named:  Named{Name: "Query"},
						Output: &Object{Type: graphql.String},
						Handler: func() string {
							return "hello"
						},
					},
				},
			},
			[]*Mutation{},
		}
		fields, _ := resolver.getFields()
		assert.Contains(t, fields, "Query")
		assert.Len(t, fields["Query"].Args, 0)
		assert.Equal(t, graphql.String, fields["Query"].Type)
	})
	t.Run("query with input", func(t *testing.T) {
		resolver := Resolver{
			Named: Named{Name: "Input"},
			Queries: []*Query{
				{
					Operation: Operation{
						Named: Named{Name: "Input"},
						Inputs: []*Input{
							{
								Object: Object{
									Named: Named{Name: "name"},
									Type:  graphql.String,
								},
							},
						},
						Output: &Object{Type: graphql.String},
						Handler: func(args struct{ name string }) string {
							return fmt.Sprintf("hello %s", args.name)
						},
					},
				},
			},
		}
		fields, _ := resolver.getFields()
		assert.Contains(t, fields, "Input")
		assert.Len(t, fields["Input"].Args, 1)
		assert.Contains(t, fields["Input"].Args, "name")
		assert.Equal(t, fields["Input"].Args["name"].Type, graphql.String)
	})
	t.Run("mutation", func(t *testing.T) {
		resolver := Resolver{
			Mutations: []*Mutation{
				{
					Operation: Operation{
						Named: Named{Name: "Mutation"},
						Inputs: []*Input{
							{
								Object: Object{
									Named: Named{Name: "name"},
									Type:  graphql.String,
								},
							},
						},
						Output: &Object{Type: graphql.String},
						Handler: func(args struct{ name string }) string {
							return fmt.Sprintf("hello %s", args.name)
						},
					},
				},
			},
		}
		_, fields := resolver.getFields()
		assert.Contains(t, fields, "Mutation")
		assert.Contains(t, fields["Mutation"].Args, "name")
		assert.Equal(t, graphql.String, fields["Mutation"].Args["name"].Type)
	})
}
