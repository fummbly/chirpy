# Chirpy

A twitter clone for posting chirps and browse other users chirps.

## Why Chirpy

I created chirpy because I wanted to practice using http server patterns. I was able to learn new techniques for server endpoints and how to properly use authentication for a server. I learned the proper ways to setup a RESP Api. I also learned more ways of using postgres and updating databases.

### Goal

I learned how to setup server endpoints and use the go package [net/http](https://pkg.go.dev/net/http). I learned how to use JWT with the package [JWT](https://pkg.go.dev/github.com/golang-jwt/jwt/v5) and dealing with refresh tokens. I used proper REST Api endpoint creation and webhook endpoints. 

## Installation

This project is not meant to be installed but if you would like to install it just clone this repo
```
git clone https://github.com/fummbly/chirpy
```
You will also need to setup [postgres](https://www.postgresql.org/) and [goose](https://github.com/pressly/goose) in order to have the backend working

## Quick Start


### Users

Users are created with a Post request to api/users.

Request body:
```
{
  "email": "user@email.com",
  "password": "password"
}
```

Sample Response body:
```
{
  "id":"8caa4762-9b4e-4c4b-8d83-12d326bc9f22",
  "created_at":"2024-12-09T10:42:51.07676Z",
  "updated_at":"2024-12-09T10:42:51.07676Z",
  "email":"email@email.com",
  "is_chirpy_red":false
}
```

For a user to login use the api/login endpoint.

Request body:
```
{
  "email": "user@email.com",
  "password": "password"
}
```

Sample Response body:
```
{
  "id": "8caa4762-9b4e-4c4b-8d83-12d326bc9f22",
  "created_at": "2024-12-09T10:42:51.07676Z",
  "updated_at": "2024-12-09T10:42:51.07676Z",
  "email": "email@email.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiOGNhYTQ3NjItOWI0ZS00YzRiLThkODMtMTJkMzI2YmM5ZjIyIiwiZXhwIjoxNzMzNzYyOTE1LCJpYXQiOjE3MzM3NTkzMTV9.uVjGwug63d_BVmTatdL91bsfNmap1Q6s7FxxfiYgLMY",
  "refresh_token":"bec4a3df7607eadf976cbf30f97db475781b4147fbc5954af9672e72f62e7c72"
}
```
The token key is the users JWT token for session login and the refresh_token is for a refresh of the session token.

## Chirps

Chirps are created with a Post request to the api/chirps endpoint.


> **_NOTE:_** A user must be logged in to create chirps.


Request body:
```
{
  "body": "chirp body text"
}
```

Request headers: 
```
{
  "Authorization": "Bearer {user JWT token}"
}
```

Sample Response body:
```
{
  "id": "{chirp_id}"
  "created_at": "{chirp creation time}"
  "updated_at": "{chirp updated time}"
  "user_id": "{authors id}"
  "body": "{body of chirp}"
}
```

Chirps can be requested by sending a GET request at the api/chirps endpoint.

Chrips can be filtered with the author query eg: ?author={authorID}

Chirps are automatically sorted in ascending order but can be changed with the sort query eg: ?sort=desc

Sample Response body:
```
[
  {"id":"24bd2a72-aa5e-455e-b037-9274152f3d2f","created_at":"0000-01-01T19:22:41.937395Z","updated_at":"0000-01-01T19:22:41.937395Z","user_id":"eb254a55-0a1e-4df7-8329-eef4938f0e68","body":"I'm the one who knocks!"},
  {"id":"3bf9d808-c58e-4387-97e4-74817fff1ebe","created_at":"0000-01-01T19:22:41.939924Z","updated_at":"0000-01-01T19:22:41.939924Z","user_id":"eb254a55-0a1e-4df7-8329-eef4938f0e68","body":"Gale!"},
  {"id":"897b78a2-b413-4c40-8755-86406ed6a0fe","created_at":"0000-01-01T19:22:41.942346Z","updated_at":"0000-01-01T19:22:41.942346Z","user_id":"eb254a55-0a1e-4df7-8329-eef4938f0e68","body":"Cmon Pinkman"},
  {"id":"f8f73bd4-7052-4e22-aa5c-8379d165ab1a","created_at":"0000-01-01T19:22:41.944676Z","updated_at":"0000-01-01T19:22:41.944676Z","user_id":"eb254a55-0a1e-4df7-8329-eef4938f0e68","body":"Darn that fly, I just wanna cook"}
]
```

Alternatively a single chirp can be requested by sending a GET request to api/chirps/{chirp_id}


## Auth 

A JWT token can be refreshed by sending a POST request at the endpoint api/refresh with the refresh_token in the header.

Request headers:
```
{
  "Authorization": "Bearer {user refresh_token}"
}
```

Sample Response body:
```
{
  "token": ""eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiOGNhYTQ3NjItOWI0ZS00YzRiLThkODMtMTJkMzI2YmM5ZjIyIiwiZXhwIjoxNzMzNzYyOTE1LCJpYXQiOjE3MzM3NTkzMTV9.uVjGwug63d_BVmTatdL91bsfNmap1Q6s7FxxfiYgLMY"
}
```

A users refresh token can be revoked with a POST request to the api/revoke endpoint with the users refresh_token in the headers.

Request headers:
```
{
  "Authorization": "Bearer {user refresh_token}"
}
```

Response:
```
204 - No Content
```
