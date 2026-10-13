# Gator RSS Feed

This is a RSS(really simple syndication) feed project. Copy paste the urls of content websites you want to get updates to build your feed yourself. This is only a backend project that runs on your cli.

## Tech Stack
- **Backend** : Go
- **Database** : Postgres
- **Migration and queries** : Goose, sqlc

## Features
Cli tool

## Installation
1. Install go compiler and postgres driver
2. run:
```terminal
go install github.com/user14556454/project2
```
3. Make a json file named .gatorconfig.json in your home:
```
{
    "db_url": <database_connection_string>,
    "current_user_name" : ""
}
```

## Set up database
1. Install postgres
2. Install goose:
```terminal
go install github.com/pressly/goose/v3/cmd/goose@latest
```
3. Go to the project root and run:
```
goose postgres <your_database_connection_string> up
```
