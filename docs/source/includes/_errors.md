# HTTP Status Codes

Ape returns the following HTTP status codes


Status Code | Meaning
---------- | -------
200 | OK -- Good Job!
201 | Created -- A new resource was created using a POST request
204 | No Content -- The server processed your request, but is not returning any content
400 | Bad Request -- Your request sucks
403 | Forbidden -- The resource requested is hidden for administrators only
404 | Not Found -- The specified resource could not be found
405 | Method Not Allowed -- You tried to access a resource with an invalid method
406 | Not Acceptable -- You requested a format that isn't json
409 | Conflict -- You tried creating a resource that already exists
429 | Too Many Requests -- You're requesting too many resources! Slow down!
500 | Internal Server Error -- We had a problem with our server. Try again later.
503 | Service Unavailable -- We're temporarially offline for maintanance. Please try again later.
