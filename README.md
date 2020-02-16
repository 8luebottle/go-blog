# go-blog
> Learn By Doing : Go Blog  

#### RESTful API with Go  
<img width="400" alt="go-blog" src="https://user-images.githubusercontent.com/48475824/74310185-e6803980-4daf-11ea-97f1-96a053d9dc5f.png">

### Table of Contents
* [Tech Stack](#tech-stack)
* [Run go-blog](#run-go-blog)
* [DB](#db)
* [Directory Structure](#directory-structure)

## Tech Stack
Go, JWT, IntelliJ, Bcrypt, GORM, Docker, MySQL, Kubernetes, Gorilla Mux

<img width="536" alt="Tech Stack" src="https://user-images.githubusercontent.com/48475824/74404146-fc9fff80-4e6c-11ea-90f9-9fce0c92a77d.png">


## Run go-blog
```bash
go run main.go
```
Successfully Connected to the MySQL
<img width="847" alt="go-blog mysql" src="https://user-images.githubusercontent.com/48475824/74598295-6eb55600-50b2-11ea-8322-534c55af5d5f.png">


## DB
Tables
```mysql
+------------------+
| Tables_in_goblog |
+------------------+
| posts            |
| users            |
+------------------+
```
Posts
```mysql
+------------+-----------------+------+-----+-------------------+-------------------+
| Field      | Type            | Null | Key | Default           | Extra             |
+------------+-----------------+------+-----+-------------------+-------------------+
| id         | bigint unsigned | NO   | PRI | NULL              | auto_increment    |
| title      | varchar(255)    | NO   | UNI | NULL              |                   |
| content    | varchar(500)    | NO   |     | NULL              |                   |
| author_id  | int unsigned    | NO   |     | NULL              |                   |
| created_at | datetime        | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
| updated_at | datetime        | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+------------+-----------------+------+-----+-------------------+-------------------+

+----+-------------+------------------------------+-----------+---------------------+---------------------+
| id | title       | content                      | author_id | created_at          | updated_at          |
+----+-------------+------------------------------+-----------+---------------------+---------------------+
|  1 | First Post  | Go Programming. Daily Coding |         1 | 2020-02-16 11:39:58 | 2020-02-16 11:39:58 |
|  2 | Second Post | Level Up Coding.             |         2 | 2020-02-16 11:39:58 | 2020-02-16 11:39:58 |
+----+-------------+------------------------------+-----------+---------------------+---------------------+
```
users
```mysql
+------------+--------------+------+-----+-------------------+-------------------+
| Field      | Type         | Null | Key | Default           | Extra             |
+------------+--------------+------+-----+-------------------+-------------------+
| id         | int unsigned | NO   | PRI | NULL              | auto_increment    |
| nickname   | varchar(255) | NO   | UNI | NULL              |                   |
| email      | varchar(150) | NO   |     | NULL              |                   |
| password   | varchar(100) | NO   |     | NULL              |                   |
| created_at | datetime     | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
| updated_at | datetime     | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+------------+--------------+------+-----+-------------------+-------------------+

+----+-------------+----------------------+--------------------------------------------------------------+---------------------+---------------------+
| id | nickname    | email                | password                                                     | created_at          | updated_at          |
+----+-------------+----------------------+--------------------------------------------------------------+---------------------+---------------------+
|  1 | Baby Tiger  | babytiger@gmal.com   | $PWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWP | 2020-02-16 11:39:58 | 2020-02-16 11:39:58 |
|  2 | Thom Browne | thombrowne@gmail.com | $PWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWPWP | 2020-02-16 11:39:58 | 2020-02-16 11:39:58 |
+----+-------------+----------------------+--------------------------------------------------------------+---------------------+------------------
```

## Directory Structure
```bash
.
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
│   │   ├── post.go
│   │   ├── user.go
│   │   └── user_test.go
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
├── tests
│   └── user_test.go
└── vendor
    ├── github.com
    │   ├── badoux
    │   │   └── checkmail
    │   │       ├── LICENSE
    │   │       ├── README.md
    │   │       └── checkmail.go
    │   └── jinzhu
    │       ├── gorm
    │       │   ├── License
    │       │   ├── README.md
    │       │   ├── association.go
    │       │   ├── callback.go
    │       │   ├── callback_create.go
    │       │   ├── callback_delete.go
    │       │   ├── callback_query.go
    │       │   ├── callback_query_preload.go
    │       │   ├── callback_row_query.go
    │       │   ├── callback_save.go
    │       │   ├── callback_update.go
    │       │   ├── dialect.go
    │       │   ├── dialect_common.go
    │       │   ├── dialect_mysql.go
    │       │   ├── dialect_postgres.go
    │       │   ├── dialect_sqlite3.go
    │       │   ├── docker-compose.yml
    │       │   ├── errors.go
    │       │   ├── field.go
    │       │   ├── go.mod
    │       │   ├── go.sum
    │       │   ├── interface.go
    │       │   ├── join_table_handler.go
    │       │   ├── logger.go
    │       │   ├── main.go
    │       │   ├── model.go
    │       │   ├── model_struct.go
    │       │   ├── naming.go
    │       │   ├── scope.go
    │       │   ├── search.go
    │       │   ├── test_all.sh
    │       │   ├── utils.go
    │       │   └── wercker.yml
    │       └── inflection
    │           ├── LICENSE
    │           ├── README.md
    │           ├── go.mod
    │           ├── inflections.go
    │           └── wercker.yml
    ├── golang.org
    │   └── x
    │       └── crypto
    │           ├── AUTHORS
    │           ├── CONTRIBUTORS
    │           ├── LICENSE
    │           ├── PATENTS
    │           ├── bcrypt
    │           │   ├── base64.go
    │           │   └── bcrypt.go
    │           └── blowfish
    │               ├── block.go
    │               ├── cipher.go
    │               └── const.go
    └── modules.txt
```