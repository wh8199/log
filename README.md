# Introduce
This a simple and easy using logging tools. Mainly used as a modular log
, The log of each module can have a name, and it will print the call time and the number of lines called, which can help developers quickly find the wrong place. Currently it supports five levels of logging output.

# Usage

## import this package
go get github.com/wh8199/log

## simple print
If you have a module named 'test', you may initialize a logging instance like below
```
logging := log.NewLogging("test", log.INFO_LEVEL, 2)
logging.Info("This is a test logging message")
```
It will print 
```
[ test ] 2020-11-08 11:40:53,332 /home/wh8199/golang/src/log-demo/main.go:11 Info msg: This is a test logging message 
```

## format
Also, if you want to format output message, you can use function with f. Next is an example
```
logging.Infof("This is a test %s message","logging")
```
It will print
```
[ test ] 2020-11-08 11:40:53,332 /home/wh8199/golang/src/log-demo/main.go:11 Info msg: This is a test logging message 
```

## logging level
When the logging instance is creating, you can select the level of log output, the built-in log level are
```
DEBUG_LEVEL   
INFO_LEVEL
WARN_LEVEL
ERROR_LEVEL
FATAL_LEVEL
```
## custom output formatter
If you don't like the default output formatter, you can custom the output format by yourself with the help of 'NewLoggingWithFormater' when you initializing logging instance

# Test and benchmark

## Test 
make test, it will print the test result

## Benchmark
make bench, it will print benchmark result

