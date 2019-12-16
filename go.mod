module github.com/jopicornell/go-rest-api

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.3.3
	github.com/bxcodec/faker/v3 v3.1.0
	github.com/disintegration/imageorient v0.0.0-20180920195336-8147d86e83ec // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/gbrlsnchs/jwt/v3 v3.0.0-rc.1
	github.com/go-jet/jet v0.0.0-00010101000000-000000000000
	github.com/go-playground/validator/v10 v10.0.1
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang-migrate/migrate/v4 v4.7.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.1
	github.com/gorilla/websocket v1.4.1
	github.com/graux/image-manager v0.0.0-20191204164210-f8171224d867
	github.com/jackc/pgx v3.2.0+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
	golang.org/x/crypto v0.0.0-20190927123631-a832865fa7ad
	golang.org/x/image v0.0.0-20191206065243-da761ea9ff43 // indirect
)

replace github.com/go-jet/jet => github.com/jopicornell/jet v2.1.3+incompatible
