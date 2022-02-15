module github.com/cosmos-gaminghub/exploder-graphql

go 1.16

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/cosmos/cosmos-sdk v0.42.4 // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/cors v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/vektah/gqlparser/v2 v2.1.0
	go.mongodb.org/mongo-driver v1.5.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
