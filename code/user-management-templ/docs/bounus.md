

now we need to tackle a big problem with the mps generation

we want to separate the Entity (user.go) files Struct Definition and the Methods in different files

so we have to use a workaround i "invented"

# plan
## Specialized Comments

```go

// PROJECT:ProjectName
// BEGIN_FILE:fileName.go
.. some code
.. some code
.. some code
// BEGIN_FILE:secondFileName.go
.. some code for secondFile

```
## Specialized Scripts and Automation
- then we need to write scripts to:
    1. find the files and read them
    2. parse the simple splitter comments and split the files
    3. make the new files in the targeted path and maintain the directory structure
- common stuff
    1. run the language specific Initializations and such ie: for golang 
        a. go mod init <ProjectName>
        b. go mod tidy // sync the deps and download them if any
        c. optionally go vet, go build etc..
    2. Dockerize the app by generating a standard DockerFile just like how we split the files by adding Specialized Comments we can extract ProjectName, Version of the App, Version of Go etc  to add to the DockerFile
    3. then if the app requires multiple apps we need to also make a docker-compose.yml





