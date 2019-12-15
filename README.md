# Images API v1.0.0
This is the 

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

## Start the project
This project uses the new go modules. This means you should have installed the latest version of go. This brings the
benefit (or drawback, if you are an old go developer 😜) that you could clone this repo in any folder of your system and
run the `go mod download` and you should be done with dependencies.

You should install and configure postgres to be able to run the migrations. After that run 
`go build -o ./bin ./cmd/migrate` to be able to run the migrations and run `chmod +x ./bin/migrate`. Then to run the
migrations is as easy as run `./bin/migrate up`.

Once migrations are ran, you can build the main app with `go build -o ./bin ./cmd/images-api`, give exec perms with 
`chmod +x ./bin/images_api` and run `./bin/images_api`

## Folder Structure
Go language doesn't define a folder structure, but there's a non-official structure that I wanted to follow, as it is
a mix between what I learnt on PalmaActiva and what I saw on different projects like docker and web boilerplates. 

### Cmd
Here are all mains for all the commands. We tried to have the mains as simple as they could be. The migrate command is
one of that commands that should be present on all projects that have persistent data on databases. All others are app
related commands like users, appointments and so on.

### Db
Here are all db related files like migrations. 

#### Migrations
Find here all the ups and downs of all migrations in sql format.

#### Entities
Here you'll find [go-jet](https://github.com/go-jet/jet) generated files like models and tables.

To generate this files again with your structure, just run:
 ```
 jet -source=PostgreSQL -host=localhost -port=5432 -user={user} -password={pasword} -dbname={db} -schema={schema} -path=./db/entities
```

### Api
The `api.go` file is the bootstrap of the app. It contains the little logic to start the handlers and the server.
Here I separate in folders the different modules this app has. You'll find some files like `handler.go` that contain all
the handler logic and `x_errors.go` that contain common used errors over the module and the app.

#### middlewares
All the middlewares that are required for this or other modules. For example, the user middleware is used all over the app,
but it is placed here as it belongs to the user universe

#### Models
I only use models when is strictly necessary. May seem like **responses** and **requests** are models as well, but they
are only used to format the input and output of the app. If a struct has to have some logic, then is treated as model.
In this app, the only model you'll see is the User, placed in the user module.

#### Requests
I've researched that coupling models to requests in go leads to problems. This is something Laravel does,
and I really like it. Here you'll find all the requests that are handled by handler, with validations.

#### Responses
The same as requests, but for the responses. Usually Jet allows you to redefine models and make very complex trees in
no time. That's where all

#### Services
A folder containing all the logic of the app. Usually, you'll see no coupling between Handlers and 

### Pkg
Here you'll find all package related source codes. This is the folder that could be transformed into libraries to be
used on other projects. Here you'll find Server wrappers, Router wrappers, utilities... So they are done in a way that 
I find productive and easy. This allows me to be able to replace external libraries with ease because they are wrapped.