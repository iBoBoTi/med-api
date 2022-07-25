
## Installation Guidelines

### Provision app using docker
To run this app using docker, follow the command below 

```bash
  git clone git@github.com:decagonhq/meddle-api.git
  cd meddle-api
  docker compose up --build
```
### Provision app locally with pg-sql

Step 1: Install and configure pgAdmin locally.

Step 2: Create a new database.

Step 3: Ask someone in the team for the .env credentials.

Step 3: Run the command below.
```bash
  go run main.go
```
### Api documentation link

```http://localhost:8080/swagger
```