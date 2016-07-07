# Formosa

Formosa is cluster database using leveldb storage data.

#Warning

##This version not final version. All futures still in developing.

# Version

version: 0.0.1

# Futures
	support functions:
		auth
		set
		get
		del
		incr
		exists
		scan
		batchexec
		hset
		hget
		hdel
		hincr
		hexists
		hsize
		hscan
		zip
		qpush
		qpop
		qsize


# Configuration

using json format to configuration.

## Configuration Example

```
	{
	  "debug":true,
	  "host":"127.0.0.1", //listen ip
	  "port":4001,// listen port
	  "password":"", //database password
	  "timeout":120, // client idle timeout
	  "mode":"mirror",//mirror or slave
	  "nodelist":[ //cluster nodes
	    {
	      "id":"db1",
	      "host":"127.0.0.1",
	      "port":4002,
	      "password":"ssdbpassword",
	      "mode":"mirror"//all  command will auto sync up to this database.
	    },
	    {
	      "id":"db2",
	      "host":"127.0.0.1",
	      "port":4003,
	      "password":"ssdbpassword",
	      "mode":"slave"//slave db
	    }
	    ]
	}
```
#How to build

```
 go get github.com/matishsiao/formosa/
 cd $GOPATH/github.com/matishsiao/formosa
 go build
```

#see more information?

#License
N/A

#Reference

https://github.com/ideawu/ssdb

https://github.com/syndtr/goleveldb
