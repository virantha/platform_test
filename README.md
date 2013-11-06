platform_test
=============
Author: Virantha Ekanayake

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

- JSON not matching schema
- Blank message
- Empty recipients list
- Recipients list larger than 5000
- Any recipient that is not 10 digits

Source
------
The source is written in Go:

[Github](https://github.com/virantha/platform_test/)

Go source files:

- [Main routine](https://github.com/virantha/platform_test/blob/master/go/src/github.com/virantha/server.go)
- [Error utility function](https://github.com/virantha/platform_test/blob/master/go/src/github.com/virantha/errors.go)
- [A unit test for the main computation](https://github.com/virantha/platform_test/blob/master/go/src/github.com/virantha/server_test.go)


Answers to questions
--------------------
a. Complexity

Each request is processed in O(n) time where n is the number of recipients. If
the category count (m) starts becoming large, then it will be O(nm)

The linear time comes from having to traverse through each recipient to validate a proper phone number.
Plus, of course, there's the underlying JSON parsing and assembling of the routed
JSON, that is also linear time.

Yes, it's possible to optimally solve this problem in polynomial time.  With
other throughput values, it may not be possible.  For instance, if you did not
have the single message relay, then you could not optimally route any message
that was not a multiple of 5.  In general, you would need your recipient count to
be an integer multiple of at least one of the throughput values.





