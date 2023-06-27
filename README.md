#Shortly Shortly is a URL shortener API that allows users to create and manage
short URLs for their own use. The API provides user authentication
functionality, as well as routes for creating, retrieving, and redirecting short
URLs.

Getting Started To get started with Shortly, you'll need to have Go and Fiber
installed on your machine. You can then clone the Shortly repository and run the
application using the following commands:

go Once the application is running, you can access the API at
http://localhost:3000/api/v1.

Authentication Shortly uses JWT-based authentication to secure its routes. To
authenticate a user, you'll need to send a POST request to the /api/v1/login
route with a JSON payload containing the user's email and password. If the
credentials are valid, the API will return a JWT token that can be used to
authenticate subsequent requests.

To authenticate a request, you'll need to include the JWT token in the
Authorization header of the request. The token should be prefixed with the
string "Bearer ", like so:

Routes Shortly provides the following routes:

POST /api/v1/login Authenticates a user and returns a JWT token.

POST /api/v1/logout Logs out the current user and invalidates their JWT token.

PUT /api/v1/users/:id Updates the user with the specified ID.

DELETE /api/v1/users/:id Deletes the user with the specified ID.

GET /api/v1/r/:redirect Redirects the user to the original URL associated with a
short URL.

GET /api/v1/shortly Retrieves all short URLs associated with the current user.

GET /api/v1/shortly/:id Retrieves the short URL with the specified ID.

POST /api/v1/shortly Creates a new short URL.

Contributing If you'd like to contribute to Shortly, please fork the repository
and submit a pull request. We welcome contributions of all kinds, including bug
fixes, new features, and documentation improvements.

License Shortly is licensed under the MIT License. See LICENSE for more
information.
