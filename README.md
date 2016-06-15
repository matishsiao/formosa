# Formosa

Formosa is cluster database using leveldb storage data.

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
		keys
		rkeys
		scan
		rscan
		batchexec
		hset
		hget
		hdel
		hincr
		hexists
		hsize
		hkeys
		hgetall
		hscan
		hrscan
		hclear
		zip 
	

# Configuration

use json format to configuration proxy setting.

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
comming soon

