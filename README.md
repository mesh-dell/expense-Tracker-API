# Expense Tracker API

A RESTful Expense Tracker API built with **Go**, **GORM**, and **JWT**, following **Clean Architecture** principles.

---

## Features

* User registration & login
* JWT access & refresh tokens
* CRUD expenses
* Expense filtering by date
* MySQL + GORM
* Clean architecture

---

## Tech Stack

* Go
* MySQL
* GORM
* JWT

---

## Project Structure

```
cmd/server        # App entrypoint
internal/         # Domain, services, repositories, handlers
config/           # App configuration
```

---

## Environment Setup

Create a `.env` file in the project root:

```env
DB_USER=root
DB_NAME=expense
DB_PASSWORD=
DB_ADDR=localhost:3306

PORT=8080

ACCESS_SECRET=
REFRESH_SECRET=
```


---

## Getting Started

```bash
git clone https://github.com/mesh-dell/expense-Tracker-API.git
cd expense-Tracker-API
go mod tidy
```

Create database:

```sql
CREATE DATABASE expense;
```

Run server:

```bash
go run cmd/server/main.go
```

Server runs on `http://localhost:8080`

---

## Authentication

* `POST /auth/register`
* `POST /auth/login`

---

## Expense Endpoints

* `POST /expenses`
* `GET /expenses`
* `PUT /expenses/{id}`
* `DELETE /expenses/{id}`

---

## Filtering

```http
GET /expenses?filter=week
GET /expenses?filter=month
GET /expenses?filter=3months
GET /expenses?filter=custom&start=startdate&end=enddate
```

---
https://roadmap.sh/projects/expense-tracker-api
