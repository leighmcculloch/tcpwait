# wait

Wait for TCP ports to be open.

Useful for starting services that need their dependencies to be available first.

## Install

```
go install 4d63.com/wait@latest
```

## Usage

```
wait -a 8080 -a postgres:5432 -t 1m
```
