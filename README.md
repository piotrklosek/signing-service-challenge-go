# Signature Service - Coding Challenge

## Challenge description

This project was created for a recruitment process.  
The task instructions can be found [here](./challenge_requirements.md).

## Instructions

### Run app

To run app we can use different approach

```
// run using default values
go run . server

// run using ENVS
SIG_PORT=8080 SIG_DB_TYPE=memory go run . server

// run using config file
go run . server --config=./config_filepath.json
```

Viper configuration is set to look for `config.json` in `./` path if config flag is not specify


All of available `ENV` can be found in [Config](./internal/config/config.go)

---
### Usage of app

To get all available command please use `go run . help`
```
Signature Service - signing API this microservice was created as a part of recruitment process

Usage:
  signature-service [command]

Available Commands:
  help        Help about any command
  migrate     Run database migrations (Postgres only)
  server      Run the Signature Service API server

Flags:
      --config string   Config file (yaml/json)
  -h, --help            help for signature-service

Use "signature-service [command] --help" for more information about a command.
```
---
### Future Improvements

The current structure was designed to be easily extendable. Some possible areas of development include:

- **Database Integrations** – prepared packages for connecting to databases such as MongoDB and PostgreSQL, with a flexible configuration layer that allows plugging in different drivers.  
- **User Management** – a set of pre-defined handlers and models to manage user entities, ready to be expanded with authentication, roles, and permissions.  
- **Session Middleware** – a planned middleware layer for handling user sessions and access control, ensuring secure and scalable session management across the application.  

These building blocks make the project a good foundation for evolving into a more complete service-oriented application.






