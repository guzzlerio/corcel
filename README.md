# Corcel

Website up at https://corcel.io

Develop branch : [![build status](http://ci.guzzler.io/ci/projects/1/status.png?ref=develop)](http://ci.guzzler.io/ci/projects/1?ref=develop)
Release branch : [![build status](http://ci.guzzler.io/ci/projects/1/status.png?ref=release)](http://ci.guzzler.io/ci/projects/1?ref=release)
Master branch   : [![build status](http://ci.guzzler.io/ci/projects/1/status.png?ref=master)](http://ci.guzzler.io/ci/projects/1?ref=master)

## What is it and what will it be?

This is a performance testing tool.  Some of the core aims of this project are:

 - Support many protocols
 - Cross platform (one simple statically compiled binary per platform)
 - Effective system resource utilization
 - Fast
 - Scalable
 - Consistent
 - Cater for all types of technical staff
 - Detailed and clear reporting
 - Support for Scenarios

## What tools already exist?

Some of the tools which exist today include:

 - HP Load Runner
 - Apache JMeter
 - Gattling
 - Apache AB
 - Siege

Each differ from the number of features, the number of supported protocols and other various things. 

## Why will this one be different?

I can't say yet but I am going to use the core aims at the top of this README and user feedback to keep me on track.

## Example

This example will use *enanos* which is another tool under guzzlerio which is simply a multi functional test http server.  To install:

```shell
go get github.com/guzzlerio/enanos
```

Next, start the enanos server as a background process or in another shell.

```shell
enanos -p 5000
```

Next, create a list of urls which you want to use.  Each line will be addressed to `http://127.0.0.1:5000/` which is the *enanos* server.  

```shell
echo http://127.0.0.1:5000/success > my-urls-to-test.txt
echo http://127.0.0.1:5000/server_error >> my-urls-to-test.txt
echo http://127.0.0.1:5000/server_error >> my-urls-to-test.txt
echo http://127.0.0.1:5000/server_error >> my-urls-to-test.txt
echo http://127.0.0.1:5000/success >> my-urls-to-test.txt
echo http://127.0.0.1:5000/success >> my-urls-to-test.txt
echo http://127.0.0.1:5000/success >> my-urls-to-test.txt
```

You can see there is a mixture of `success` and `server_error` in the list, which is how you can use *enanos* as it has specific endpoints which return the different ranges of http response codes.  The above will produce 200 and a random selection of the 5XX response range.

Next, make sure the source is built and yu have the *code-named-something* executable and then invoke with the following arguments:

```shell
./code-named-something -f ./my-urls-to-test.txt --summary --workers 10
```

Once it has finished you will then see console output similar to the following:

```shell
Running Time: 0.007 s
Throughput: 9381 req/s
Total Requests: 70
Number of Errors: 30
Availability: 57.14285714285714%
Bytes Sent: 3510
Bytes Received: 8990
Mean Response Time: 0.1429 ms
Min Response Time: 0 ms
Max Response Time: 1 ms
```

```
_
```
