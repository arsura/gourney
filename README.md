Gourney.
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
- [mockery](https://github.com/vektra/mockery) -  Tool for generating mock files for Golang interfaces [Download](https://github.com/vektra/mockery/releases)
- [ifacemaker](https://github.com/vburenin/ifacemaker) - Tool for generating Golang interfaces
- [migrate](https://github.com/golang-migrate/migrate) - Tool for migration Database [Download](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

Commands
-----

### Generate MongoDB's Collection Interface (mongo.Collection)
```sh
ifacemaker -f  ~/go/pkg/mod/go.mongodb.org/mongo-driver@v1.9.1/mongo/collection.go -s Collection -i MongoCollectionProvider -p adapter -o .type.go
```

### Generate MongoDB's Client Interface (mongo.Client)
```sh
ifacemaker -f  ~/go/pkg/mod/go.mongodb.org/mongo-driver@v1.9.1/mongo/client.go -s Client -i MongoClientProvider -p adapter -o type.go
```

### Generate MongoDB's Database Interface (mongo.Database)
```sh
ifacemaker -f  ~/go/pkg/mod/go.mongodb.org/mongo-driver@v1.9.1/mongo/database.go -s Database -i MongoDatabaseProvider -p adapter -o type.go
```

### Generate Mock Files
```sh
mockery --all --keeptree --dir ./pkg/adapters 
```

File Structure
-----

```bash
├── cmd
│   ├── api
│   ├── worker
│   └── main.go
├── config
│   ├── config.go
│   └── config.yml
├── deployments
│   └── Dockerfile.api
├── developments
│   ├── docker-compose.infra.yml
│   ├── docker-compose.yml
│   └── Dockerfile.api
├── pkg
│   ├── adapters
│   ├── constant
│   ├── logger
│   ├── models
│   ├── repositories
│   ├── services
│   ├── usecases
│   ├── utils
│   └── validator
 ```