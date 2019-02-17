## Cronsched

Schedules a crontab-like (not real crontab) list of commands to run.

### Install, build, test

To install, clone this repo into your Go path.

The application has a single third-party dependency. It uses the cheekbits/is
test package to make writing tests in Go a little easier.

```
$ go get github.com/cheekybits/is
```

To build a `cronsched` binary (from `cmd/main.go`) in the root directory, run:

```
$ make build
```

To run tests, run:

```
$ make test
```

### Usage

Pass the command list in via stdin.

From a file:

```
$ cat fixture/crontab | ./cronsched
```

Using hereoc:

```
$ ./cronsched <<EOF               
45 * /bin/run_me_hourly
* * /bin/run_me_every_minute
* 19 /bin/run_me_sixty_times
EOF
```

The output of the above run at 19:02 will produce:

```
19:45 today - /bin/run_me_hourly
19:02 today - /bin/run_me_every_minute
19:00 tomorrow - /bin/run_me_sixty_times
```

Time can be manually passed in via `-time` flag:

```
$ ./cronsched -time 19:02
```

