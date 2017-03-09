# cloud collector

## why do this

I need to collection some system info on many servers:
-  how many connections some app created to mongo
-  how many connections established by mongo self(different form above)
 
 if there is no limits, mongodb may go down because of too many connetions.

## how to  use 

pull this repo, follow what I have done

```
{
    "cmds": ["./shell/mongostat.sh", "./shell/linux_netstat.sh"],
    "metrics": ["mongo_tatal_", "conn_mongo_"],
    "interval": 60,
    "endpoint": "www"
}
```

In this config, I set to shell cmd to collect data, the data will be parsed and then pushed to [open-falcon](https://github.com/open-falcon), which has dashboard to show the metrics and alert me if need

## has question?

email me, I will reply to you