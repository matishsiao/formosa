# Formosa

Formosa is cluster database using leveldb storage data.

#Warning

##This version not final version. All features still in developing.

# Version

version: 0.0.1

# Features

## support functions:

		auth
		set
		get
		del
		incr
		exists
		scan
		rscan
		size
		getall
		batchexec (using json format)
		batchwrite (using json format)
		hset
		hget
		hdel
		hincr
		hexists
		hsize
		hscan
		hrscan
		hgetall
		qpush
		qpop
		qsize
		zip

## websocket

you can use websocket to subscribe all changed event in real time.

```
	 var ws = new WebSocket("ws://127.0.0.1:8080/sub",["hashname","keyname"]);
		 ws.onmessage = function(e) {
			 console.log("Frames Receive:" + event.data);
		 };

 ```

# Functions
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

Command:rscan

	reverse scan KV data in range.

Example:

```
//limit:-1 = no limit
client.Do("rscan","from","to",limit)
```

Command:getall

	get all in key value map.

Example:

```
client.Do("getall")
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
client.Do("hscan","hashname","from","to",limit)
```

Command:hrscan

	reverse scan hash KV data in range.

Example:

```
//limit:-1 = no limit
client.Do("hrscan","hashname","from","to",limit)
```

Command:hgetall

	get all in one hash map.

Example:

```
client.Do("hgetall")
```

Command:hsize

	hash KV data size

Example:

```
client.Do("hsize","hash")
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

Command:zip

	using gzip for transfer data. Save bandwidth resource.

Example:

```
//turn ON
client.Do("zip",1)
//turn OFF
client.Do("zip",0)
```


# Configuration

using json format to configuration.

## Configuration Example

```
	{
	  "debug":false,
	  "host":"127.0.0.1",							//db listen ip
	  "port":4001,										//db listen port
	  "dbpath":"data",								//db storage folder
	  "web":{ 												//web service setting
	    "host":"127.0.0.1",
	    "port":8080
	  },
	  "password":"defaultpassword",
	  "timeout":120, 									//unit:second
	  "limit":5, 											//Node connection limit
	  "nodelist":[
	    {
	      "id":"db2",
	      "host":"127.0.0.1",
	      "port":4002,
	      "password":"defaultpassword"
	    },
	    {
	      "id":"db3",
	      "host":"127.0.0.1",
	      "port":4003,
	      "password":"defaultpassword"
	    }
	  ]
	}
```
#How to build

```
 go get github.com/matishsiao/goformosa/
 cd $GOPATH/github.com/matishsiao/goformosa
 go build server.go
```

#See more information?

[![GoDoc](https://godoc.org/github.com/matishsiao/goformosa/formosa?status.svg)](https://godoc.org/github.com/matishsiao/goformosa/formosa)

#License

Copyright 2016 Matis Hsiao <matismaya@gmail.com> All rights reserved.

#Reference

https://github.com/ideawu/ssdb

https://github.com/syndtr/goleveldb
