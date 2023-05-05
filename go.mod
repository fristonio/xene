module github.com/fristonio/xene

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/dgraph-io/badger v1.6.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.9.0
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/golang/protobuf v1.5.0
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20190717161051-705d9623b7c1+incompatible
