# Formosa

Formosa is cluster database using leveldb storage data.

#Warning

##This version not final version. All features still in developing.

# Version

version: 0.0.1

# Features
	support functions:
		auth
		set
		get
		del
		incr
		exists
		scan
		size
		batchexec (using json format)
		batchwrite (using json format)
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

# functions
Command:auth

	use password to verify client

Example:

```
	client.Do("auth","defaultpassword")
```

Command:set

	set key|value data to formosa

Example:

```
	client.Do("set","key","value")
```

Command:get

	get key|value data

Example:

```
	client.Do("get","key")
```

Command:del

	delete key|value data

Example:

```
	client.Do("del","key")
```

Command:incr

	atomic operation to incr value data, It will return incr value.

Example:

```
	client.Do("incr","key",1)
```

Command:exists

	check key exist in formosa.

Example:

```
	client.Do("exists","key")
```

Command:scan

	scan KV data in range.

Example:

```
  //limit:-1 = no limit
	client.Do("scan","from","to",limit)
```

Command:size

	KV data size

Example:

```
	client.Do("size")
```

Command:batchexec

	use batch to run commands in one command. this data must be json string.

Return:

	result array

Example:

```
	client.Do("batchexec","[[\"hset\",\"test\",\"1\",\"1\"],[\"hget\",\"test\",\"1\"]]")
```

Command:batchwrite

	use batch to write data in one command. this data must be json string.


Example:

```
	client.Do("batchwrite","[[\"hset\",\"test\",\"1\",\"1\"],[\"set\",\"KV\",\"1\"]]")
```

Command:hset

	write hash key data to formosa.


Example:

```
	client.Do("hset","hash","key","value")
```

Command:hget

	get hash key data from formosa.


Example:

```
	client.Do("hget","hash","key")
```

Command:hdel

	delete hash key data from formosa.


Example:

```
	client.Do("hdel","hash","key")
```

Command:hincr

	atomic operation to incr value data, It will return incr value.


Example:

```
	client.Do("hincr","hash","key",1)
```

Command:hexists

	check key exist in formosa.

Example:

```
	client.Do("hexists","hash","key")
```

Command:hscan

	scan hash KV data in range.

Example:

```
  //limit:-1 = no limit
	client.Do("hscan","hash","from","to",limit)
```

Command:hsize

	hash KV data size

Example:

```
	client.Do("hsize","hash")
```

Command:zip

	using gzip for transfer data. Save bandwidth resource.

Example:

```
	//turn ON
	client.Do("zip",1)
	//turn OFF
	client.Do("zip",0)
```

Command:qpush

	write data to formosa queue.


Example:

```
	client.Do("qpush","queue","value")
```

Command:qpop

	get queue first element from formosa.


Example:

```
	client.Do("qpop","queue")
```

Command:qsize

	queue data size

Example:

```
	client.Do("qsize","queue")
```



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

[![GoDoc](https://godoc.org/github.com/matishsiao/goformosa/formosa?status.svg)](https://godoc.org/github.com/matishsiao/goformosa/formosa)

#License

Copyright 2016 Matis Hsiao <matismaya@gmail.com> All rights reserved.

#Reference

https://github.com/ideawu/ssdb

https://github.com/syndtr/goleveldb
