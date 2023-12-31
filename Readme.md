# MARKETSPACE API 🚀

![Capa](/images/Capa.png)

Api made in Golang for the marketspace

technologies used:
  - Go 1.2
  - Go-chi v1.5.4
  - Go-chi/jwtauth v5.1.0
  - Swaggo/http-swagger v1.3.4
  - Swaggo/swag v1.16.1
  - Docker v20.10.17
  - Docker-compose v2.10.2

## How to Run?

### Docker and Docker-compose

if you have docker and docker-compose, you can run this project with just one command:

```docker-compose
  docker-compose up -d
```

this command will upload the database and create the necessary tables and will build the api. The database init script is inside db_scripts folder if you want to modify it.

<hr/>

### Go Run

if you don't want to install docker-compose just use docker to upload the database:

```docker
  # change the container-name and passoword
  docker run --name container-name -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres

```

Just remember to change env variables inside .env file.

To run the api just use the command:

```golang
  go run main.go -env=

```
* env
  * development
  * production

env flag is optional, if you don't pass it then .env is loaded otherwise you can pass production then .env.production file is loaded.

<hr/>

### Swagger

To access swagger, you just access this [http://localhost:8000/v1/swagger/index.html#/]("http://localhost:8000/v1/swagger/index.html#/").
