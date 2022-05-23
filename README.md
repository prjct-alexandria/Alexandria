# Alexandria
We work for clients Andrew Demetriou, PhD. student, and Cynthia Liem, PhD., to develop a collaborative platform and workflow for scientific publishing, for the duration of quarter 4. Our goal is to produce an open-source collaborative scientific publishing software product by translating the git workflow that is usually utilised by a narrow demographic (programmers, computer scientists) into a process that scientists from all domains will be able to use to publish articles and perform peer reviewing.

# Usage
## API Documentation
An API documentation can be generated from go annotations by [swag](https://github.com/swaggo/swag). Run from the project root:
    
    cd ./mainServer
    swag init -g server/router.go

The documentation can be updated with the same commands, when the endpoints change in the code.

After generating it, when running the backend, it can be accessed through http://localhost:8080/docs/index.html. (Or the equivalent when running it on a different server address.)



# Team
| Name                | Email adress                       |
|---------------------|------------------------------------|
| Amy van der Meijden | A.vanderMeijden@student.tudelft.nl |
| Andreea Zlei        | A.Zlei@student.tudelft.nl          | 
| Jos Sloof           | A.J.G.Sloof@student.tudelft.nl     |
| Mattheo de Wit      | M.C.A.deWit@student.tudelft.nl     |
| Emiel Witting       | E.A.Witting@student.tudelft.nl     |



***
