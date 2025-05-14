# go-api

## Description

A simple starter template for building a future, more robust Go-based API.
It includes the following features:

* Logging: Records request data.
* Authentication: Verifies requests using JWTs.
* MongoDB Integration: Connects to a MongoDB database.

## Environment Variables

### 1. MongoDB Connection

Configure your MongoDB connection using the following variables:

```env
MONGO_URI       = MongoDB connection string  
DB_COLLECTION   = Name of the collection  
DB              = Name of the database  
```

### 2. JWT Creation

These variables are required to generate JWTs:

```env
COMPANY         = Name of the company  
COMPANY_EMAIL   = Contact email address  
```

## Modes

The application supports two modes:

1. **API Server (default)**

   ```
   go run .
   ```

2. **JWT Generator**

   ```
   go run . -mode=jwt
   ```

## Routes

### 1. `/status`

* **Method**: `GET`
* **Purpose**: Returns the HTTP status of the API

---

### 2. `/enter`

* **Method**: `POST`
* **Purpose**: Submit an entry
* **Required Fields**:

  * `name`
  * `email`
  * `age`
  * `token`

---

### 3. `/getByEmail`

* **Method**: `POST`
* **Purpose**: Retrieve user data by email
* **Required Fields**:

  * `email`
  * `token`

---


## Example Request

```
curl -X POST http://127.0.0.1:8080/enter \                        
     -H "Content-Type: application/json" \
     -d '{"name":"Adrian","email":"arojo@arojo.com", "age":25, "token":"TestToken"}' 
```
