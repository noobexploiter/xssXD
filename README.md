# xssXD

A tool made to detect xss vulnerablities in a list of urls. It takes input from the stdin.

# How to install 
```go get github.com/noobexploiter/xssXD```

# How to use
```
Usage of ./xssXD:
  -c int
        Set the Concurrency  (default 50)
  -s string
        Specify the payload to use (default "none")
```
* Set Concurrency according to your need, default 50
* Specify the payload to use. If not specified, it will check for characters <'"> by default

`cat urls.txt | xssXD -c 100`

# Additional Info

The list of urls must be in the format ```protocol://subdomain/path?querys```

# Example
## Default
```
cat test | ./xssXD                                                                                                   
http://public-firing-range.appspot.com/reflected/parameter/body/400?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/title?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/attribute_singlequoted?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/body/500?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/body_comment?q=a is reflecting <, ', ", >
https://xss-game.appspot.com/level1/frame?query=asd is reflecting <, ', ", >
http://sudo.co.il/xss/level0.php?email=asd# is reflecting <, ', ", >
```
## With specified payload
```
cat test | ./xssXD -s "<svg/onload=alert()>"                                                                         
http://public-firing-range.appspot.com/reflected/parameter/title?q=a is reflecting <svg/onload=alert()>
http://public-firing-range.appspot.com/reflected/parameter/body/400?q=a is reflecting <svg/onload=alert()>
http://public-firing-range.appspot.com/reflected/parameter/body/500?q=a is reflecting <svg/onload=alert()>
http://public-firing-range.appspot.com/reflected/parameter/body_comment?q=a is reflecting <svg/onload=alert()>
http://public-firing-range.appspot.com/reflected/parameter/attribute_singlequoted?q=a is reflecting <svg/onload=alert()>
https://xss-game.appspot.com/level1/frame?query=asd is reflecting <svg/onload=alert()>
http://sudo.co.il/xss/level0.php?email=asd# is reflecting <svg/onload=alert()>
```
## Set Concurrency
```
cat test | ./xssXD -c 100
http://public-firing-range.appspot.com/reflected/parameter/attribute_singlequoted?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/body/400?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/body/500?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/body_comment?q=a is reflecting <, ', ", >
http://public-firing-range.appspot.com/reflected/parameter/title?q=a is reflecting <, ', ", >
https://xss-game.appspot.com/level1/frame?query=asd is reflecting <, ', ", >
http://sudo.co.il/xss/level0.php?email=asd# is reflecting <, ', ", >
```

# UPDATES
## 11/3/2020
I added ssti payload {{7\*7}} to detect ssti too.
## 11/6/2020
Added -v option for verbose mode
