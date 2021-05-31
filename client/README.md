# FYST Data Migration Pipeline Client 
**Author:** *Ivan Sharankov*


## Overview
This is the first (unofficial) version of the FYST Data Migration Pipeline, a pipeline intended to log, categorize, transport, and copy all the data generated at the Fred Young Submillimeter Observatory.


## Installation
### Installing Golang

To use this package, you will need to install Golang (Go) to be able to build the binaries to run the Simulator. The binaries for Golang can be installed from official Golang page found [here](https://golang.org/dl/). Also, [this link](https://golang.org/doc/install/source) may be useful for installing Go from source. It is important to correctly setup your `$GOPATH`, which can be explained for multiple operating systems [here](https://github.com/golang/go/wiki/SettingGOPATH). If Golang is installed correctly, you can use `go env` to see the Go environment information. Note: if you choose to setup Golang in a directory that is not the default, you will need to [set your $GOROOT](https://www.tecmint.com/install-go-in-linux/) respectively.

Once Golang is set up, you can learn the Golang basics from [here](https://tour.golang.org/welcome/1). Generally golang.org has been designed in such a way that you can learn absolutely everything you need to start programming from just their domain, meaning it's an invaluable tool for both new and experienced Gophers (Golang programmers). 


### Running the Client
To use the client, you will need to build the package with `go install`.



___
## Interface 
### Client Commands
The client provides a simple way to communicate with the FYST Data Migration Pipeline with a single command. Currently three endpoints exist, as follows:


#### Upload
The upload endpoint can be used by passing the keyword `upload` after the binary. This endpoint has several flags that can be used to populate metadata associated with the data being uploaded. Two of these flags are required for a successful upload, which are `md5` and `name`. The last four (`instrument`, `date`, `size`, `url`) are optional. Using this command would look something like the following.

```
client upload -name=AUniqueFilename -instrument=2 -md5=c4ca4238a0b923820dcc509a6f75849b -date=2021-05-31T12:00:00Z -size=10 -url=/path/filename.txt
```

although the `url` flag is not necessary for the server, it will be needed to locate the file on the local directory to upload it, so currently the client will fail uploading without it.


#### Rules
This endpoint simply asks the server for the rules within the database. The only flag included is `id`, which would define a single ruleID being requested.
```
client rules -id=2
```


#### Files
This endpoint asks the server for all the files within the database, and returns a maximum of the last 50 entries. The only flag included is `filename`, which would define a single Filename being requested.
```
client rules -filename=AUniqueFilename
```
