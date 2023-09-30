# My Personal Template for Go

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

~~Don't~~ use it for production!

## Clone the project

```
$ git clone https://github.com/stonear/go-template
$ cd go-template
```

## Running Template

```
$ // create .env file
$ make dep
$ go run .
```

## Create Migration

```
$ dbmate n <your-migration-name>
$ dbmate migrate
```

## Rollback Migration

```
$ dbmate down
```

## Generate SQLC

```
$ make gen
```

## Lint

```
$ make lint
```

## REST Resource Naming Guide

Verb | URI | Action
-----|-----|-------
GET | /entities | Index
~~GET~~ | ~~/entities/create~~ | ~~Create~~ (Not implemented)
POST | /entities | Store
GET | /entities/:id | Show
~~GET~~ | ~~/entities/:id/edit~~ | ~~Edit~~ (Not implemented)
PUT | /entities/:id | Update
DELETE | /entities/:id | Destroy

## Commit Naming Rules
A typical git commit message will look like: ```<type>: <subject>```, "type" must be one of the following mentioned below:
- ```build```: Build related changes
- ```chore```: A code change that external user won't see (eg: change to .gitignore file or .prettierrc file)
- ```feat```: A new feature
- ```fix```: A bug fix
- ```docs```: Documentation related changes
- ```refactor```: A code that neither fix bug nor adds a feature. (eg: You can use this when there is semantic changes like renaming a variable/ function name)
- ```perf```: A code that improves performance
- ```style```: A code that is related to styling
- ```test```: Adding new test or making changes to existing test

## Testing


```
$ make test
```

Files with the suffix ```_test.go``` are only compiled and run by the ```go test``` tool.