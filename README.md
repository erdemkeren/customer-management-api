# Customer Management System

This is a simple Customer Management System implemented in Go using the Gin web framework.


## Features

- List all customers
- Get a customer by ID
- Create a new customer
- Update an existing customer
- Delete a customer


## Installation

1. Make sure you have Go installed. You can download it from the official website: [https://golang.org](https://golang.org)

2. Clone this repository:

```bash
git clone https://github.com/erdemkeren/customer-management-api.git

cd customer-management-api

go mod download

go run main.go
```

3. The server will start running on http://localhost:8080.

4. You can use a tool like cURL or Postman to interact with the API. Postman collection is included in the repo.


## API Endpoints:

- List all customers: GET /customers
- Get a customer by ID: GET /customers/:id
- Create a new customer: POST /customers
- Update an existing customer: PUT /customers/:id
- Delete a customer: DELETE /customers/:id


## Notes

- Uses `encoding/json` under the hood.

```
https://gin-gonic.com/docs/jsoniter/
```


## License:

This project is licensed under the MIT License.
