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