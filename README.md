# Await

Wait for a file to be created/deleted/modified while blocking.

## Usage

```
await [--exists] [--notexists] [--fresh t] [--interval t] [FILE]
```

The `await` command will repeatedly check if a path exists (by default), does not exist,
or has been modified in the last `t` amount of time. The intended use case is to call `await`
as part of a shell script. For example, perhaps you're using a shell script to start a server:
you could invoke `await` to wait until the server created a particular file.

### Examples

Wait for a file to exist, as part of a shell script.

```bash
#!/bin/sh

start_server() {
    local PID_FILE=$1

    # do some thing here

    echo $! > $PID_FILE
}

wait_for_server_start() {
    start_server /run/server.pid
    await --exists /run/server.pid
}

wait_for_server_start
```

Wait for a file to be removed.

```bash
await --notexists /tmp/import-job-running.lock
```
Wait for a file to no longer be "fresh" (i.e. updated in the last 5 minutes)

```bash
await --fresh 5m /var/log/my-server.log
echo "Server isn't writing to its log!!!"
```

### Options

* `--exists` Check that the provided `FILE` exists and exit as soon as it does
* `--notexists` Check that the provided `FILE` does not exist and exit as soon as it does not
* `--fresh t` Check if that the provided `FILE` has been modified in the last `t` duration and
  exit as soon as it has not. `t` is a "duration" argument and so accepts values like `1s`, `5m`,
  `1h`, etc.
* `--interval t` How long to wait between checking that a file has been created/deleted/modified.
  `t` is a "duration" argument and so accepts values like `1s`, `5m`, `1h`, etc. The default is
  one second (`1s`).
* `FILE` File path to check for existing/not existing/freshness.

### Exit Codes

* `await` will exit with status code `0` once the given condition has been satisfied.
* `await` will exit with status code `1` if invalid input was supplied of if the file
  could not be checked for some reason.

## Building

Await is a basic Go project and doesn't require anything special to build, just the
standard library.

```
git clone git@github.com:56quarters/await.git && cd await
go build
./await --exists /tmp/foo
```
