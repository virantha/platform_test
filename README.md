platform_test
=============


Usage
-----
Please use the following URL to access the REST endpoint for "sending" a message:

    http://techvectors.com:8080/message/route

The endpoint satisfies a POST request with the specified schema. If the POST request
is valid, it will return a JSON object ("Content-Type", "application/json") in the
specified format.  Otherwise, it will return a HTTP error code along with the error
message.


Error-checking
--------------
The following conditions will trigger an error response:

- JSON not matching spechema
- Blank message
- Empty recipients list
- Recipients list larger than 5000



