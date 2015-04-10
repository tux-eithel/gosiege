# GOSIEGE, http/https stress tester made in GO
***

New in golang world, try some concurrency.

Inspired by [Siege](http://www.joedog.org/siege-home/) I have try to do something similar. 
I have added some useful features like:
* filter response header
* send custom header in request
* send custom param in POST request

### Install

```
go get -u github.com/tux-eithel/gosiege/
```

If you have installed [Godep](https://github.com/tools/godep)
```
godep restore
```

### Command Line

```
-c x
	Where _x_ is number of concurrent connections. Default 1

-exp "HeaderField Value"
	Filter header using regular expresion. You can define multiple regexp.
	*HeaderField* is a header field, *Value* must be a regular expression. 
	Examples are:
	* `-exp "X-Cache HIT"`
	* `-exp "X-Cache .*"`
```

```
-f fileName
```
*fileName* which contains for every row at least an url


```
-nasty=boolValue
```
Use all available CPUs, physic and logic. Default true


```
-per y
```	
Test every url *y* times. When all url are tested *y* time, gosiege quits.
Default -1, so test will run until you press Ctrl+c	