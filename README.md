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

2. Install dependencies:
```bash
go mod tidy
```

3. Configure your  settings in `conf/app.conf`

4. Run the application:
```bash
go run main.go
```

## API Endpoints

- `GET /v1/api/property/details/:propertyId` - List all properties
- `GET /v1/api/propertyList?propertyIds=prop-1,prop-2,prop-3 ` - Returns formated details output of matching property of available property ids
- `GET /v1/api/property/:propertyId /gallery/` - Returns labelled output of images having `confidence>95`


## Testing

Run all tests:
```bash
go test -v ./tests
```

Generate coverage report:
```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
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