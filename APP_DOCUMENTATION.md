#Documentation of the app
This document are the requeriments for the app, as well as they serve as a way to see an oversight of what the app is
capable of do.

##Locations
A location is simply a physical store, a place where the app could manage clients and customers could make appointments.
Usually a user has only one location, but in case of chains or franchises, a manager can manage multiple locations.

###Create a location [POST /locations]
A location will need a name, address, phone, start hour, end hour, appointment duration and default appointment span. 
An appointment duration is the duration of a single slot of time. This is used to measure how many appointments a location
can have. An appointment span is how many slots are set as default for an appointment.

Additionally, it can have an email and a list of employeeId in case it is not the first location we create and we have 
already a list of employees. This endpoint needs the user to be a *manager*.

###Update a location [PATCH /locations]
You can update the name, the address, the phone and the list of employees. This endpoint needs the user to be a *manager*.

###Delete a location [DELETE /locations]
You can delete any location if it has no appointments.This endpoint needs the user to be a *manager*.

###List all employees of a location [GET /locations/{id}/users]
This endpoint needs the user to be a *manager*.

###Create a new employee to a location [POST /locations/{id}/users]
This endpoint needs the user to be a *manager*.
###Delete an employee of a location [DELETE /locations/{id}/users/{id}]
This endpoint needs the user to be a *manager*.

##Availability and Time Slots [GET /slots?available=1&duration=30&date=20190930]
To check for availability, slots are checked. Each location have the appointment duration field, which indicates
how many minutes a slot have. For example, from 9am to 10am, with an appointment duration of 15min, there are 4 slots.

##Appointments
An appointment is a time slot on a calendar of a location and employee. It has a start time and a duration in minutes.


##Users
Any user can have any role. This means that a user can be an employee and a manager at the same time.
###Roles
An **employee** is the type of user in charge of the main use of the app: taking appointments and manage them in place.
A **manager** is the type of user in charge of administrating the locations and employees of a location.

###Create a user [POST /users]
You'll be able to define the roles of the user here, as well as if the user has some location associated, which ones
(list of location ids).
###Update a user [PATCH /users]
Update email, name and locations
###Delete a user [DELETE /users]
You can delete a user if it's not an employee and has no future appointments

###List all appointments of a user [GET /users/{id}/appointments]
If the user is an employee, this lists all appointments of the user.

###List all locations of a user [GET /users/{id}/locations]


