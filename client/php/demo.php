<?php

include(dirname(__FILE__) .'/formosa/Formosa.php');
set_time_limit(900);
ini_set('memory_limit', '1024M');
$host = '127.0.0.1';
$port = 4001;
$pwd = "123";


originTest($host,4004,$pwd);
proxyTest($host,4001,$pwd,true);

//use gzip
function zipTest($host,$port,$pwd,$zip) {
	try{
		$Formosa = new SimpleFormosa($host, $port,120000);
	}catch(Exception $e){
		die(__LINE__ . ' ' . $e->getMessage());
	}
	$auth = $Formosa->auth($pwd);
	if ($auth == null) {
		//echo "auth success";
	} else {
		echo "auth failed";
		exit();
	}
	$data = null;
	if ($zip) {
		$result = $Formosa->zip("1");
	}
	$start = microtime(true);
	//$data = $Formosa->hscan("OneTableTest","","",2000);
	$data = $Formosa->hgetall("OneTableTest");
	$useTime = microtime(true) -$start;
	$len = count($data);
	echo "proxyTest use:$useTime len:$len\n";
}

function nonZipTest($host,$port,$pwd) {
	try{
		$Formosa = new SimpleFormosa($host, $port,120000);
		//$Formosa->easy();
	}catch(Exception $e){
		die(__LINE__ . ' ' . $e->getMessage());
	}
	$auth = $Formosa->auth($pwd);
	if ($auth == null) {
		//echo "auth success";
	} else {
		echo "auth failed";
		exit();
	}
	$start = microtime(true);
	$data = $Formosa->hgetall("OneTableTest");
	$useTime = microtime(true) - $start;
	$len = count($data);
	echo "originTest use:$useTime len:$len\n";
}
