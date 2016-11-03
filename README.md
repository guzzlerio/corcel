![Corcel](http://docs.corcel.io/images/corcel-logo.png)

# Corcel

Website up at https://corcel.io

Develop branch : [![Build Status](https://travis-ci.org/guzzlerio/corcel.svg?branch=develop)](https://travis-ci.org/guzzlerio/corcel)
Release branch : [![Build Status](https://travis-ci.org/guzzlerio/corcel.svg?branch=release)](https://travis-ci.org/guzzlerio/corcel)
Master branch   : [![Build Status](https://travis-ci.org/guzzlerio/corcel.svg?branch=master)](https://travis-ci.org/guzzlerio/corcel)

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
 - httperf
 - wrk
 - vegeta
 - autocannon

and more ...

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
./corcel run --summary --workers 5 ./my-urls-to-test.txt 
```

Once it has finished you will then see console output similar to the following:

```shell
╔═══════════════════════════════════════════════════════════════════╗
║                           Summary                                 ║
╠═══════════════════════════════════════════════════════════════════╣
║         Running Time: 10.001014785s                               ║
║           Throughput: 2349 req/s                                  ║
║       Total Requests: 23487                                       ║
║     Number of Errors: 0                                           ║
║         Availability: 100.0000%                                   ║
║           Bytes Sent: 1.9 MB                                      ║
║       Bytes Received: 3.0 MB                                      ║
║   Mean Response Time: 0.4669 ms                                   ║
║    Min Response Time: 0.0000 ms                                   ║
║    Max Response Time: 21.0000 ms                                  ║
╚═══════════════════════════════════════════════════════════════════╝
```

