# GolangStudying

sudo netstat -tulpn | grep :12300
sudo kill -9 3581

FreePort.sh <номер порта>


# Golang microservices assignment

## Time limits

The evaluation result of the test is not linked to how much time you spend on it.
Please DO NOT spend more than 2 hours doing it.
Successful applications show us that 1 hour is more than enough to cover all the evaluation points below.

This assignment is meant to evaluate the golang proficiency of full-time engineers.
Your code structure should follow microservices best practices and our evaluation will focus primarily on your ability to follow good design principles and less on correctness and completeness of algorithms. During the face to face interview you will have the opportunity to explain your design choices and provide justifications for the eventually omitted parts.

**Do not mention 90 percent of everything or 90poe anywhere on the code or repository name.**

**Do not copy any files from this repository**

## Evaluation points in order of importance

- use of clean code which is self documenting
- use of packages to achieve separation of concerns
- use of domain driven design
- use of golang idiomatic principles
- use of docker
- tests for business logic
- use of code quality checkers such as linters and build tools
- use of git with appropriate commit messages
- documentation: README and inline code comments

Results: please share a git repository with us containing your implementation.

Level of experience targeted: EXPERT

Avoid using frameworks such as go-kit and go-micro since one of the purposes of the assignment is to evaluate the candidate ability of structuring the solution in their own way.
If you have questions please make some assumptions and collect in writing your interpretations.

Good luck.

Time limitations: time is measured from when we send you the assignment to the final commit time onto your repository.

## Technical test

- Given a file with ports data (ports.json), write 2 services.
- First service(ClientAPI) should parse the JSON file and have REST interface.
- The file is of unknown size, it can contain several millions of records.
The service has limited resources available (e.g. 200MB ram).
- While reading the file, it should call a second service(PortDomainService) that either creates a new record in a database, or updates the existing one.
- The end result should be a database containing the ports, representing the latest version found in the JSON. Database can be Map.
- First service(ClientAPI) should provide an endpoint to retrieve the data from the second service(PortDomainService)
- Each service should be built using Dockerfile
- Provide all tests that you think are needed for your assignment. This will allow the reviewer to evaluate your critical thinking as well as your knowledge about testing.
- Please use gRPC as a transport between services.
- Readme should explain how to run code and tests

Choose the approach that you think is best (ie most flexible).

## Bonus points

- Database in docker container
- Domain Driven Design
- Docker-compose file

## Note

As mentioned earlier, the services have limited resources available, and the JSON file can be several hundred megabytes (if not gigabytes) in size.
This means that you will not able to read the entire file at once.

We are looking for the ClientAPI (the service reading the JSON) to be written in a way that is easy to reuse, give or take a few customisations.
The services themselves should handle certain signals correctly (e.g. a TERM or KILL signal should result in a graceful shutdown).