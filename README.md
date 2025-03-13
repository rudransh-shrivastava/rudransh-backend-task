# Setup Process

####Clone the repo
```bash
git clone https://github.com/rudransh-shrivastava/rudransh-backend-task.git
```

####Build docker image
```bash
sudo docker build -t api-server .
```

####run docker image
```bash
docker run -p 8080:8080 api-server
```


# Project Structure Overview

├── Dockerfile           (Contains the instructions to dockerise the api server)
├── Makefile             (Define commands like make run)
├── bin                 (Static Binaries are compiled here)
│   └── server    
├── cmd                 (The main entry point of our application)
│   └── server
│       └── main.go
├── db.sqlite3          (The database file)
├── go.mod
├── go.sum
├── internal            
│   ├── api             (The api handlers, and the api server)
│   │   ├── api.go
│   │   ├── auth.go
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── middleware.go
│   │   └── mockquiz.go
│   ├── config              (Handles the environment variables if any)
│   │   └── config.go
│   ├── db                 (The configuration of our database)
│   │   └── db.go
│   ├── schema             (The schemas of our database is defined here)
│   │   └── schema.go
│   ├── store              (Store handles communication with the database)
│   │   ├── course.go
│   │   ├── quiz.go
│   │   ├── store.go
│   │   └── user.go
│   └── utils
│       ├── logger
│       │   └── logger.go   (The logger configurations)
│       └── utils.go        (Reusable utilities that we use in our code)
└── key.json                (Firebase credentials obtained from firebase)


# Database Schema and API Overview

# How to run tests?
To run the tests, please run the following command.
```bash
go test ./... -v
```
