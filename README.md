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
