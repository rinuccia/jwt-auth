# HTTP Authentication Server

This is an example implementation of jwt auth with Golang.

## Requirements

* docker & docker-compose

## Run Project

* Configure `.env` file
* Use `docker-compose up -d` to build and run docker containers with application and mongodb instance

## API Endpoints

### POST `"/auth/sign-up"`

Registers new user

Request
```json
{
  "first_name": "Maria",
  "last_name": "Petrova",
  "password": "abc123456",
  "email": "petrova@gmail.com",
  "user_type": "ADMIN"
}
```

Response
HTTP Status Code 200
```json
{
  "InsertedID": "633a0b9f7e610fed6a6f8c84"
}
```

### POST `"/auth/sign-in"`

Generates JWT Token and RefreshToken
Request
```json
{
"email": "petrova@gmail.com",
"password": "abc123456"
}
```

Response
HTTP Status Code 200
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6Im1hdHNrZXZpY2hAZ21haWwuY29tIiwiRmlyc3ROYW1lIjoiS2F0c2lhcnluYSIsIkxhc3ROYW1lIjoiTWF0c2tldmljaCIsIlVpZCI6IiIsIlVzZXJUeXBlIjoiQURNSU4iLCJleHAiOjE2NjQ4ODY4MDZ9.-T3Yb1R_kklmRzg_9JVePovKgKBk-7PHh-ADfg4Nw_Y",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0TmFtZSI6IiIsIkxhc3ROYW1lIjoiIiwiVWlkIjoiIiwiVXNlclR5cGUiOiIiLCJleHAiOjE2NjU0MDQ4ODJ9.wsmHDKxMTMgeJiicSt-JIfXRJd7aFZdKI-bdhr_1n70"
}
```