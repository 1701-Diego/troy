# Troy - Empathic status reports

Troy can fetch the status of the Tasks and LRPs for a given Domain

## To Install

```
go get github.com/1701-diego/troy
```

## To Use

First export the `RECEPTOR` address

For Ketchup:
```
export RECEPTOR=http://username:password@receptor.ketchup.cf-app.com
```

For Diego-Edge:
```
export RECEPTOR=http://receptor.192.168.11.11.xip.io
```

Then, to view the Tasks and LRPs for a given domain:

```
troy DOMAIN
```