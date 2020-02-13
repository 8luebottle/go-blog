# go-blog
> Learn By Doing : Go Blog  

#### RESTful API with Go  
<img width="500" alt="go-blog" src="https://user-images.githubusercontent.com/48475824/74310185-e6803980-4daf-11ea-97f1-96a053d9dc5f.png">

## Tech Stack
Go, JWT, IntelliJ, Bcrypt, GORM, Docker, MySQL, Kubernetes, Gorilla Mux

![tech stack](https://user-images.githubusercontent.com/48475824/74404146-fc9fff80-4e6c-11ea-90f9-9fce0c92a77d.png)

## Directory Structure
```
├── README.md
├── api
│   ├── auth
│   │   └── token.go
│   ├── controllers
│   │   ├── base.go
│   │   ├── home_controller.go
│   │   ├── login_controller.go
│   │   ├── posts_controller.go
│   │   ├── routes.go
│   │   └── users_controller.go
│   ├── middlewares
│   │   └── middlewares.go
│   ├── models
│   │   ├── Post.go
│   │   └── User.go
│   ├── responses
│   │   └── json.go
│   ├── seed
│   │   └── seeder.go
│   ├── server.go
│   └── utils
│       └── formaterror
│           └── formaterror.go
├── go.mod
├── go.sum
├── main.go
└── tests
```
