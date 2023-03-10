# Simple REST API

A simple REST API for CRUD services with 2 sample entities
# Features
## User : 
- Register new user (admin/user)
- Login user
- Get a user data (show profile)
- Show all users data (only admin)
- Edit user data (update profile)
- Soft delete

## Item :
- Insert new item
- Show one or all items
- Update item detail
- Soft delete

# API Documentation

[Click here](https://documenter.getpostman.com/view/23707537/2s8Z72WCSk) to see the documentation.

# Tools & Requirements

- Go 1.19.3
- Gorm & MySQL
- Echo v4
- JWT

## How to Install

- Clone it

```
$ git clone https://github.com/hebobibun/golang-rest-api
```


- Go to directory

```
$ cd golang-rest-api
```

- Create a new database

- Rename `local.env.example` to `local.env`
- Adjust `local.env` as your environment settings

- Run the project

```
$ go run .
```

# Enjoy

Keep learning! ^^
