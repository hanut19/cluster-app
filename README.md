# Cluster Management Portal

A simple web application built in Go that allows users to view and manage cluster node counts. The application supports two user roles: `admin` (can view and update clusters) and `readonly` (can only view clusters). User authentication is session-based, and PostgreSQL is used as the backend database.

## Features

- User login/logout with role-based access
- View a portal of clusters and their node counts
- Admins can update node counts via the UI
- Middleware to enforce login and admin access
- HTML templates for login and portal views

## Project Structure

 ``` cluster-app/ 
     ├── cmd/ # Application entry point
     │    └── main.go
     ├── db/ # Database connection and schema
     │    ├── db.go
     │    └── schema.sql
     ├── handlers/ # HTTP handlers (Login, Logout, Portal, Update)
     │    └── handlers.go
     ├── middleware/ # Session and role-based middleware
     │    └── session.go
     ├── templates/ # HTML templates
     │    ├── login.html
     │    └── update.html
     ├── config.json # Database configuration file (not included by default)
     └── README.md ```


## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL

### Step 1: Set Up the Database

Run the following SQL using `psql` or a GUI tool:

```bash
psql -U youruser -d yourdb -f db/schema.sql
