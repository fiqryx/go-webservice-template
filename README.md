## ⚡Go Webservice API Template
minimal starter template for building Go REST APIs using:

- ⚙️ [Gin](https://gin-gonic.com/) – Fast and flexible HTTP web framework
- 🧬 [GORM](https://gorm.io/) – ORM for PostgreSQL and more
- 🐘 [PostgreSQL](https://www.postgresql.org/) – Reliable and powerful relational database

### How to Run
#### 1. Clone the Repository

```bash
git clone https://github.com/fiqryx/go-webservice-template.git
cd go-webservice-template
```

#### 2. Setup Environment Variables
Copy the example .env file and configure it as needed:
```bash
cp .env.example .env
```

#### 3. Run the Application
Start http server with default configuration:
```bash
go run . serve
# OR
go run . serve --host=127.0.0.1 -p 9000
```

#### 4. (Optional) Enable Hot Reloading with [Air](https://github.com/air-verse/air)
Install Air if you don't have it:
```bash
go install github.com/air-verse/air@latest
```
Run the server with hot reload
```bash
air serve
# or
air serve --host=127.0.0.1 --port=9000
```

---

### Database Migration
Run database migrations using GORM's AutoMigrate feature. This will automatically create or update tables based on your Go model definitions.

```bash
go run . migrate
# or
go run . migrate -D
```
Make sure to register your models in the `DBRegistry.models` at [/registry/database.go](registry/database.go).

### Database backup
command to backup database with registry tables:
```bash
go run . db:backup
# or
go run . db:backup --output=./storage/backup/20250518
```
Make sure to register your tables `DBRegistry.tables` at [/registry/database.go](registry/database.go).

---

### Additional CLI Commands
This project also supports other command-line operations:

#### Create repository:
```bash
go run . make:repo --name=user
```
After created adjust registry at [/registry/repository.go](registry/repository.go).

#### Create service:
```bash
go run . make:service --name=auth
```
After created adjust registry at [/registry/services.go](registry/services.go).

#### Create controller:
```bash
go run . make:controller --name=home
```
After created adjust registry at [/registry/controller.go](registry/controller.go).

#### Create factory
command to create factory:
```bash
go run . make:factory --name=user
```

with specific output directory, default directory is [/database/factory](database/factory/)
```bash
go run . make:factory --name=user --output=./factory
# or
go run . make:factory --n user -o ./factory
```

#### Database seeding
command to run database seed with the `factories`:
```bash
go run . db:seed
```
Make sure the configuration `DBRegistry.factories` at [/registry/database.go](registry/database.go).

---

### How to use [module.sh](moduel.sh)
This script helps rename the Go module path in [go.mod](go.mod) and across your project files.
```bash
# using default values (e.g., template.go → webservices)
./module.sh

# using custom module paths
./module.sh "old/module/path" "new/module/path"
```

### How to Build
```bash
$ go build -o ./bin/api
# or with vendor
$ go build -mod=vendor -o ./bin/api
```