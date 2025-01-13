# Property Fetch Format API

A RESTful API built with Beego framework for fetching and formatting property data.

## Overview

This project provides an API service to fetch property information and format it according to specific requirements. Built using the Beego framework in Go.

## Prerequisites

- Go 1.16 or higher
- Beego framework

## Installation

1. Clone the repository:
```bash
git clone https://github.com/uzzalcse/property-fetch-format-api.git
cd property-fetch-format-api
```

2. Configure your  settings in `conf/app.conf`

3. Build the app and install dependencies:
```bash
docker-compose build
```

4. Run the application:
```bash
docker-compose up -d
```

## API Endpoints
### Property Endpoints
- `GET /v1/api/property/details/:propertyId` - List all properties
- `GET /v1/api/propertyList?propertyIds=prop-1,prop-2,prop-3 ` - Returns formated details output of matching property of available property ids
- `GET /v1/api/property/:propertyId /gallery/` - Returns labelled output of images having `confidence>95`

### User Endpoints 
- `POST /v1/api/user/` - Create a new user.
- `GET /v1/api/user/:identifier` - Get user details by identifier.
- `PUT /v1/api/user/:identifier` - Update user details by identifier.
- `DELETE /v1/api/user/:identifier` - Delete user by identifier.

## To interact with the apis go to the following link of `Swagger UI`

http://localhost:8080/swagger/index.html/index.html

## To check api using `Postman` use the following url if you didn't change the `port` in `app.conf` file.

http://localhost:8080/v1/api/user/

## Testing

First run application 

```
docker-compose up -d
```


Open the bash shell 

```
docker-compose exec app bash
```


Run all tests:
```bash
go test -v ./tests
```

Generate coverage report:
```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
```

View coverage in terminal:
```bash
go tool cover -func=coverage.out
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request