module github.com/apache/incubator-devlake

go 1.19

require (
	github.com/aws/aws-sdk-go v1.44.242
	github.com/cockroachdb/errors v1.9.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.9.1
	github.com/go-errors/errors v1.4.2
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-playground/validator/v10 v10.14.1
	github.com/gocarina/gocsv v0.0.0-20220707092902-b9da1f06c77e
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/jackc/pgx/v5 v5.3.1
	github.com/jmespath/go-jmespath v0.4.0
	github.com/lib/pq v1.10.2
	github.com/libgit2/git2go/v33 v33.0.6
	github.com/magiconair/properties v1.8.5
	github.com/manifoldco/promptui v0.9.0
	github.com/merico-dev/graphql v0.0.0-20221027131946-77460a1fd4cd
	github.com/mitchellh/hashstructure v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/panjf2000/ants/v2 v2.4.6
	github.com/robfig/cron/v3 v3.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.6.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.8.3
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.1
	github.com/tidwall/gjson v1.14.3
	github.com/viant/afs v1.16.0
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	go.temporal.io/api v1.7.1-0.20220223032354-6e6fe738916a
	go.temporal.io/sdk v1.14.0
	golang.org/x/crypto v0.9.0
	golang.org/x/exp v0.0.0-20221028150844-83b7d23a625f
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602
	golang.org/x/sync v0.2.0
	gorm.io/datatypes v1.0.1
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.1
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.3.0 // indirect
)

require (
	github.com/KyleBanks/depth v1.2.1
	github.com/Microsoft/go-winio v0.5.0
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7
	github.com/acomagu/bufpipe v1.0.3
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/cockroachdb/logtags v0.0.0-20211118104740-dabe8e521a4f
	github.com/cockroachdb/redact v1.1.3
	github.com/cpuguy83/go-md2man/v2 v2.0.2
	github.com/davecgh/go-spew v1.1.1
	github.com/emirpasic/gods v1.12.0
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a
	github.com/fsnotify/fsnotify v1.5.1
	github.com/getsentry/sentry-go v0.12.0
	github.com/gin-contrib/sse v0.1.0
	github.com/go-git/gcfg v1.5.0
	github.com/go-git/go-billy/v5 v5.3.1
	github.com/go-openapi/jsonpointer v0.19.6
	github.com/go-openapi/jsonreference v0.20.2
	github.com/go-openapi/spec v0.20.9
	github.com/go-openapi/swag v0.22.3
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gogo/googleapis v1.4.1
	github.com/gogo/protobuf v1.3.2
	github.com/gogo/status v1.1.0
	github.com/golang-jwt/jwt/v5 v5.0.0-rc.1
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/hcl v1.0.0
	github.com/imdario/mergo v0.3.12
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99
	github.com/jinzhu/inflection v1.0.0
	github.com/jinzhu/now v1.1.5
	github.com/josharian/intern v1.0.0
	github.com/json-iterator/go v1.1.12
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351
	github.com/kr/pretty v0.3.0
	github.com/kr/text v0.2.0
	github.com/leodido/go-urn v1.2.4
	github.com/mailru/easyjson v0.7.7
	github.com/mattn/go-colorable v0.1.11
	github.com/mattn/go-isatty v0.0.19
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/mitchellh/go-homedir v1.1.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.2
	github.com/pborman/uuid v1.2.1
	github.com/pelletier/go-toml v1.9.3
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/robfig/cron v1.2.0
	github.com/rogpeppe/go-internal v1.8.1
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/sergi/go-diff v1.1.0
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.6-0.20200504143853-81378bbcd8a1
	github.com/stretchr/objx v0.5.0
	github.com/subosito/gotenv v1.2.0
	github.com/tidwall/match v1.1.1
	github.com/tidwall/pretty v1.2.0
	github.com/ugorji/go/codec v1.2.11
	github.com/xanzy/ssh-agent v0.3.0
	go.uber.org/atomic v1.9.0
	golang.org/x/mod v0.10.0
	golang.org/x/net v0.10.0
	golang.org/x/sys v0.8.0
	golang.org/x/term v0.8.0
	golang.org/x/text v0.9.0
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac
	golang.org/x/tools v0.9.3
	google.golang.org/appengine v1.6.7
	google.golang.org/genproto v0.0.0-20220222213610-43724f9ea8cf
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/ini.v1 v1.62.0
	gopkg.in/warnings.v0 v0.1.2
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

//replace github.com/apache/incubator-devlake => ./
