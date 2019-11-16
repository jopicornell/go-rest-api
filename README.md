# v1.0.0
This app is aimed to be an integral solution for small business managers and their employees, but each module can be
used by anyonone. It is thought to be modular, fast and simple. The first module that will be included in this first 
version will be an appointment management.

## Appointments module
Basically, it is an appointment manager app. The flow is simple: 
 - The client calls the store or uses some app (google maps?) to make the appointment.
    - Manager receives the appointment (via app or call) and assigns the employee
    - Or Employee receives the unassigned appointment and makes him as the employee in charge
 - Client goes to the appointment
 - Employee marks him as waiting customer (optional step) 
 - Employee picks up the customer
 - Employee finishes the task and mark as finished the appointment

From the perspective of a manager (the company owner or an administrator), the app will allow to set up the locations,
employees, manage client data, see, set and assign the appointments. From the perspective of an employee, he will be
able to see all appointments of the day, see all unasigned appointments, assign to him one unassigned appointment and
mark as finished an appointment.

## Architecture
The architecture of the project is aimed to be a SoA/Microservices app. Each service is configured to be a separated
Go http server, allowing us to scale the app as it grows. It will be thought as an stateless server from the beginning,
allowing us to containerize each module separately.

## Start the project
This project uses the new go modules. This means you should have installed the latest version of go. This brings the
benefit (or drawback, if you are an old go developer 😜) that you could clone this repo in any folder of your system and
run the `go mod download` and you should be done with dependencies.

You should install and configure mariadb and create a schema so you can run this project without any problem. All tests
are mocked so running `go test ./...` will run tests and you don't need any db for that tests.

## Folder Structure
Go language doesn't define a folder structure, but there's a non-official structure that I wanted to follow, as it is
a mix between what I learnt on PalmaActiva and what I saw on different projects like docker and web boilerplates. 

### Cmd
Here are all mains for all the commands. We tried to have the mains as simple as they could be. The migrate command is
one of that commands that should be present on all projects that have persistent data on databases. All others are app
related commands like users, appointments and so on.

### Db
Here are all db related files like migrations and big SQL files.

#### Migrations
Find here all the ups and downs of all migrations in sql format.

### Internals
Those are internals of the app. Here you'll find all source code that is app related and could not be shared among other
projects/apps. There are a few concepts that I find interesting and are used on all modules:

#### Api
Here I separate in folders the different modules this app has. Usually it matches all the commands (except from the 
migrate command). Inside the Api folder there are two folders and a few repeated files. 

First we have the **handlers** 
and the **services**. **Handlers** are small controllers that are in charge of calling the correct services to do the
right work. Only "*manage*" the request and make them understandable to services. **Services**, on the other hand, are 
interfaces where all the business logic happens. *api.go* files are the "*main*" of the api. They prepare the http
server to be able to handle the requests appropriately. *routes.go* prepare all the routes with its "*ConfigureRoutes*"
method.

### Pkg
Here you'll find all package related source codes. This is the folder that could be transformed into libraries to be
used on other projects. Here you'll find Server wrappers, Router wrappers, utilities... So they are done in a way that 
I find productive and easy. This allows me to be able to replace external libraries with ease as they are wrapped.