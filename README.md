Development
===============

Run Infrastructure
-----
```sh
make dev-docker-infra
```

Run API Service
-----
```sh
make dev-api
```

Run Logging Worker Service
-----
```sh
make dev-worker
```

Tools 
-----
- [mockery](https://github.com/vektra/mockery) -  generates mocks for golang interfaces [Download](https://github.com/vektra/mockery/releases)
- [migrate](https://github.com/golang-migrate/migrate) - database migrations [Download](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [ifacemaker](https://github.com/vburenin/ifacemaker) - generates a Golang interface

Commands
-----

### Generate MongoDB Collection Interface (mongo.Collection)
```sh
ifacemaker -f  ~/go/pkg/mod/go.mongodb.org/mongo-driver@v1.9.1/mongo/collection.go -s Collection -i MongoCollectionProvider -p adapter -o .type.go
```

### Generate MongoDB Client Interface (mongo.Client)
```sh
ifacemaker -f  ~/go/pkg/mod/go.mongodb.org/mongo-driver@v1.9.1/mongo/client.go -s Client -i MongoClientProvider -p adapter -o type.go
```

### Generate MongoDB Database Interface (mongo.Database)
```sh
ifacemaker -f  ~/go/pkg/mod/go.mongodb.org/mongo-driver@v1.9.1/mongo/database.go -s Database -i MongoDatabaseProvider -p adapter -o type.go
```

### Generate Mock Files
```sh
mockery --all --keeptree --dir ./pkg/adapters 
```