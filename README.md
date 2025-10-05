# GoLang extension of courseapp

Claude describes the application:

* **Backend** : Go with Gorilla Mux router
* **Database** : ***Google Cloud Datastore***
* **Authentication** : Custom session management + ***Google SSO***
* **Static Files** : Serving HTML/JS files directly

The current implementation includes:

1. **Basic User Authentication**
   * Login/session management
   * Role-based permissions (but not fully utilized)
2. **Course Structure**
   * Course CRUD operations
   * Module CRUD operations
   * Elements (quiz questions/content)
   * Thread discussions
3. **Project Management**
   * Basic project structure for submissions

## Google Cloud datastore and Google+ OAuth2 SSO

```
This is a simple example of a Golang RESTful API,
using the Google Cloud Datastore and Google+ OAuth2 SSO.

Currently the OAuth2 SSO is only working for designated test users.

REST verified with Postman (http://localhost:8000/user route)

```

Using Repository Interface per Model (best DB migration practice),
thanks to suggestions from [Praveen](https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/).

## Creating a user

Use Postman or curl to create a user:

```
curl -X POST -H "Content-Type: application/json" -d '{"username":"Dave","email":"dave@gmail.com","password":"asdf","firstname":"Dave","password":"Norton"}' http://localhost:3000/user
```

The user object looks something like this:

```
{
    "username": "David",
    "password": "asdf",
    "email": "test@asdf.com",
    "firstname": "David",
    "lastname": "Norton"
}
```

## About the authentication model

Sessions are tracked in memory by uuid, expiry, and username.
Upon successful authentication (via /login or /sso route), a session
is created and stored in memory, and a cookie with session_token
is sent to the client.  Without active session, a 401 is returned.

The only route that currently depends on the session_token cookie
is the /user (GET) route (GetAllUsers).

To require authentication, add *validateSession* middleware to the route.  E.g.,

```
router.HandleFunc("/user", validateSession(http.HandlerFunc(userHandler.GetAll))).Methods("GET")
```

The login route currently uses cleartext for authentication;
just POST something like this in the body to the /login route
(or just login from the index page).

```
{
	"username": "Dave",
	"password": "asdf"
}
```

** @author Norton 2022
