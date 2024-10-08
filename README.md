# MarCom-Backend

MarCom-Frontend is a repository containing the Backend (web-server) code for the [MarCom](https://github.com/Hanyue-s-FYP) project

## Implementation & Features

- Implemented with Golang and it's powerful standard library as well as JWT and gRPC

## Setup and running the project

> Ensure that you have go > 1.22 installed in your system

### Clone this repository to local
> No need to recurse submodules as a copy of the generated grpc files are attached in the repository and will be updated together with the repository

```sh
git clone https://github.com/Hanyue-s-FYP/MarCom-Backend.git
```

### Setup environment variables

> Setup environment variables (with reference to .env.example), you can set the environment variables of the system directly or use a file (usually .env or .env.development), if you used .env or .env.development, change the NewConfig to pass in the env file, if you set the environment variables of the system directly, set it to empty string.
> For mail, try to stick with gmail, you can follow this link: [Legacy App Support](https://support.google.com/accounts/answer/185833) to obtain email and a password that can be used

### Install dependencies

```sh
go mod download
```

### Create Database

```sh
# macOS/linux
cat create_db.sql | sqlite3 marcom.db
# windows
sqlite3 marcom.db < .\create_db.sql
```

### Compile and run
> CGO_ENABLED=1 should be set when running the application

```sh
# build and run with temp binary
go run .
```
Alternatively,
```sh
# build a binary and store it to the build folder
go build -o build/MarCom_WebServer .
# Execute the built binary
MarCom_WebServer
```
