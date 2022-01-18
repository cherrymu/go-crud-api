# Go API
This is a repo with a simple API developed using Go performing basic CRUD operations when connected to a postgresql database using basic HTTP methods which can be deployed on kubernetes platform.

This application will expose below endpoints for HTTP calls 

**GET** - /movies/ 

**POST** - /movies/

**PUT** -  /movies/{movieid}

**DELETE** -/movies/{movieid}

**DELETE** - /movies/


To **run** and test the code locally, use below instructions

Step 1: Create a local postgresql database and load the scheme file init.sql from the repo 

Step 2: Uncomment and pass your local DB config in the below section in index.go file and comment the var section below it so that we can use this for local testing 

```
const (
	DB_USER     = "db_username"
	DB_PASSWORD = "db_password"
	DB_NAME     = "db_name"
	DB_HOST     = "db_servername"
	DB_PORT     = db_port_number
  )
```

comment below block of code for local testing using a local postgresql database

```
var (
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
)
```
Step 3: Run the code using below command 

```
go run index.go 
 ```
 
Step 4: Now, our Go app will be listening on the port 8001 and use localhost to fetch the data from the database using the route 

```
curl localhost:8001/movies/ 
```

### To Deploy this code on Kuberenetes environment please use the Helm Charts instructions mentioned in the below repo 

[Helm Chart for go-api](https://github.com/cherrymu/go-app-helm-chart)










