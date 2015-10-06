[![Stories in Ready](https://badge.waffle.io/guzzlerio/code-named-something.svg?label=ready&title=Ready)](http://waffle.io/guzzlerio/code-named-something)
[![Stories in Progress](https://badge.waffle.io/guzzlerio/code-named-something.svg?label=in progress&title=Progress)](http://waffle.io/guzzlerio/code-named-something)
[![Build Status](https://travis-ci.org/guzzlerio/code-named-something.svg?branch=develop)](https://travis-ci.org/guzzlerio/code-named-something)

# Code-named-something

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
./code-named-something -f ./my-urls-to-test.txt --summary --workers 100
```

Once it has finished you will then see console output similar to the following:

```shell
Running Time: 0.049 s
Throughput: 14020 req/s
Total Requests: 700
Number of Errors: 300
Availability: 57.14285714285714%
Bytes Sent: 35100
Bytes Received: 89900
Mean Response Time: 2.157 ms
Min Response Time: 0 ms
Max Response Time: 17 ms
```

## Git branching and release strategy

 - All issues to be worked on inside a feature file
 - All issues completed to merged into the develop branch
	- Binaries are generated for all platforms as a pre-release set of artefacts under the label `latest`
 - Upon a code freeze (when a release candidate is to be made) the `HEAD` of the `develop` branch will be merged into the `release` branch.
	- Binaries are generated for all platforms as a pre-release set of artefacts under the label `pre-release`
 - Following any hot-fixes to the release candidate the `HEAD` of the `release` branch will be merged into the `master` branch.
	- The repository will be tagged at this stage with the next version for the release artefacts.  (Need to confirm the order and possibly update the CI release build script to ensure it is sync'd)

# Progress
[![Throughput Graph](https://graphs.waffle.io/guzzlerio/code-named-something/throughput.svg)](https://waffle.io/guzzlerio/code-named-something/metrics)
