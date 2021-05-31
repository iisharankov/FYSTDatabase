# FYST Data Migration Pipeline 
**Author:** *Ivan Sharankov*


## Overview
This is the first (unofficial) version of the FYST Data Migration Pipeline, a pipeline intended to log, categorize, transport, and copy all the data generated at the Fred Young Submillimeter Observatory.


## Installation
### Installing Golang

To use this package, you will need to install Golang (Go) to be able to build the binaries to run the Simulator. The binaries for Golang can be installed from official Golang page found [here](https://golang.org/dl/). Also, [this link](https://golang.org/doc/install/source) may be useful for installing Go from source. It is important to correctly setup your `$GOPATH`, which can be explained for multiple operating systems [here](https://github.com/golang/go/wiki/SettingGOPATH). If Golang is installed correctly, you can use `go env` to see the Go environment information. Note: if you choose to setup Golang in a directory that is not the default, you will need to [set your $GOROOT](https://www.tecmint.com/install-go-in-linux/) respectively.

Once Golang is set up, you can learn the Golang basics from [here](https://tour.golang.org/welcome/1). Generally golang.org has been designed in such a way that you can learn absolutely everything you need to start programming from just their domain, meaning it's an invaluable tool for both new and experienced Gophers (Golang programmers). 


### Running docker-compose
A docker-compose YAML file is located at the base of the repository that can be started. It will automatically create a MySQL database, build the Data Migration Pipeline schema, and populate the required rows with some example rows. The YAML file also builds two instances of Minio object storage, and the volumes required to store persistent data. 

The server is currently commented out within the docker container, and must be start manually after the database is initialize. This is simply due to constant development, and as long as a new version of the /DataMigration/Dockerfile is built, the container can be uncommented within the YAML file.






___
## Interface 
### Client Commands
The simplest way to communicate with the FYST Data Migration Pipeline is by using the client provided. This can be found under the README in the client directory.

___
### Curl Commands
If for some reason the client is not preferable, HTTP requests can be sent directly to the server.  Currently three endpoints are defined and accessible, shows below: 

#### Files
The files endpoint is used to access any data that is passing through the FYST Data Migration Pipeline. There exist two HTTP GET requests for this endpoint, which are 
```
curl -GET 'localhost:8700/files'
curl -GET 'localhost:8700/files/FileName'
```
The first of these returns the last 50 files in the database in descending order. Functionality to modify this number may be added. The second request will return the metadata for a given `FileName` provided. All communication between the user and data on the server is defined by the unique `FileName` provided when the file was added to the database, and serves as a unique ID.

To add a file to the database, a HTTP POST request is sent, which should take the form:
```
curl 'localhost:8700/files' -d@- <<___
{
    "name": "NameOfTheFile",
    "instrument": 1,
    "md5sum": "c4ca4238a0b923820dcc509a6f75849b",
    "date_created": "2021-05-31T12:00:00Z",
    "size": 1024302,
    "url": "URL"
}
___
```
The URL field may be depreciated in the future. Currently both the `name` and `md5sum` fields must be given for a  `/files` POST request to be successful, the rest of the rows are optional.


To add a record to the copies endpoint, a HTTP POST request is sent to `files/FileName/copies`, which looks like
```
curl 'localhost:8700/files/FileName/copies' -d@- <<___
{
    "location_id": 1,
}
___
```
The location can be found in the body of the reply from the `\files` response.


#### Records

Just like with the `\files` endpoint, you can query a single FileName or many filenames with very similar syntax:

```
curl -GET 'localhost:8700/records'
curl -GET 'localhost:8700/records/FileName'
```
Again, the maximum number of values returned by the first request is 50. Name of the endpoint is pending....

The last `/records`/ endpoint is the `/records/FileName/LocationID` endpoint, which is a POST request that adds a record to the database. It does not contain any body.

```
curl -GET 'localhost:8700/records/FileName/LocationID'
```
Again, the locationID is a value given by the `/files` endpoint when adding a file.
#### Rules

The rules endpoint provides all the logic for how data is moved around when added to the database. Rules can be enabled or disabled to change the behaviour, though this isn't yet added to the front end. Below are the two POST requests the `/rues` endpoint accepts, with similar behaviour as the last examples:

```
curl -GET 'localhost:8700/rules'
curl -GET 'localhost:8700/rules/RuleID'
```