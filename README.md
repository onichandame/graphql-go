# Graphql Go

golang graphql framework inspired by type-graphql.

# Get Started

```go
import (
    "encoding/json"

    "github.com/gin-gonic/gin"
    graphqlGo"github.com/onichandame/graphql-go"
    graphql "github.com/graphql-go/graphql"
)

type User struct{
    Name string `graphqlgo:"name=name,type=string"`
}

var users = make([]*User,0)

type UserInput struct{
    Input *User `graphqlgo:"name=input"`
}

func (u*User)MutationcreateUser(ctx *gin.Context, args *UserInput)*User{
    users=append(users,args.Input)
    return users[len(users)-1]
}

func main(){
    schema:=graphqlGo.GetSchemaFromStruct(&User{})
    http.HandleFunc("/graphql",func(w http.ResponseWriter,r*http.Request){
        query:=r.URL.Query().Get("query")
        result:=graphql.Do(graphql.Params{
            Schema: schema,
            RequestString:query,
        })
        json.NewEncoder(w).Encode(result)
    })
    http.ListenAndServe(":80",nil)
}
```
