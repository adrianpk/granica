# Package auth

## Draft notes

### Auth

Defines methods and functions that allow to attach features provided by this package to [service module](https://gitlab.com/mikrowezel/backend/service/).


### Server

Router configurationi: endpoint functions are associated with a url (resource).

### Endpoint

Action handlers entrypoints, HTTP and/or gRPC data is decoded into transport objects, sent to service and then output received from it is encoded in order to be send to the client.
Lightweight input validations and warnings output messages must be processed in this layer.

### service

All business logic shoud reside here.

### transport

Business objects used as input/output interfaces between endpoint and service layers.

### helper

Basic helpers for the whole package. Currently only methods that allow transferring values from tranport interfaces to business objects and the other way round.

### xxx_test

Package tests.

## Links
[Home](/README.md)
