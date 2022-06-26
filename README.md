# Alexandria

---
We work for clients Andrew Demetriou, PhD. student, and Cynthia Liem, PhD., to develop a collaborative platform and workflow for scientific publishing, for the duration of quarter 4. Our goal is to produce an open-source collaborative scientific publishing software product by translating the git workflow that is usually utilised by a narrow demographic (programmers, computer scientists) into a process that scientists from all domains will be able to use to publish articles and perform peer reviewing.

# Getting started

---
This section explains usage of the software with the default configuration. If a configuration file is edited, values like ports and passwords may be different.

## Using docker compose
All the different components can be started together using docker-compose:

    docker-compose up
After running this command from the project root, the UI will be accessible at `localhost:3000`. No extra installations apart from docker and docker-compose are required.

Note that this option will create git repositories for internal use within the Alexandria repository, which might be annoying during the development process.
This option uses the `dockerconfig.json` file together with `Dockerfile` to correctly configure all ports and URL's. Changes to `mainServer/config.json` will be ignored using this option.

## Without docker
### Database
A postgreSQL database should be running on port 5432, with user `postgres` and (by default) the password `admin`.
A database should be created with the name `AlexandriaPG`. When running the system later, the program will automatically add the necessary tables.

### Back-end
To start up the backend server, have Go version 1.18 installed. Then, the project root, run:
    
    cd ./mainServer 
    go build mainServer
    go run mainServer

### Front-end
To start up the frontend server, with npm installed, from the project root, run:

    cd ./mainClient
    npm install
    npm run start-win

If not on Windows, run `npm run start` instead.

## API Documentation
A [SwaggerUI](https://swagger.io/tools/swagger-ui/) API documentation is available by default at `http://localhost:8080/docs/index.html`. If the configuration has been changed, the first part of the URL should be that of the backend server.


# Contributing

---
For contribution to the project and setting up a development environment, see [CONTRIBUTING.md](CONTRIBUTING.md)


# Team

---
| Name                | Email adress                       |
|---------------------|------------------------------------|
| Amy van der Meijden | A.vanderMeijden@student.tudelft.nl |
| Andreea Zlei        | A.Zlei@student.tudelft.nl          | 
| Jos Sloof           | A.J.G.Sloof@student.tudelft.nl     |
| Mattheo de Wit      | M.C.A.deWit@student.tudelft.nl     |
| Emiel Witting       | E.A.Witting@student.tudelft.nl     |

