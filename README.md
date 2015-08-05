# envconf

It's simple env based config file transform for docker boot time configuration.


```
$cat a.conf
[servera]
HOST={{ .HOST | default "B" }}

$export HOST="HOSTA"

$envconf a.conf

$cat a.conf
```

``` 
$envconf a.conf b.conf ...
or
$export ENV_CONF_FILES=a.conf:b.conf:c.conf
$envconf
```





