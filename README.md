# My Personal Template for Go

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

~~Don't~~ use it for production!

## Clone the project

```
$ git clone https://github.com/stonear/go-template
$ cd go-template
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

## Testing

If you are using VSCode, I recommend you to use [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) Extension for running ```apitest.http```.

To check if the database connected successfully, you can run the following command:

```go test -timeout 30s -run ^TestDb$ github.com/stonear/go-template/database```

Files with the suffix ```_test.go``` are only compiled and run by the ```go test``` tool.