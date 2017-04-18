[![Build Status](https://travis-ci.org/kayac/Gunfish.svg?branch=master)](https://travis-ci.org/kayac/Gunfish)

# Gunfish
APNS provider server on HTTP/2.

* Gunfish provides the only interface as the APNS provider server.

## Overview
![overviews 1](https://cloud.githubusercontent.com/assets/13774847/14844813/17035232-0c95-11e6-8307-1d8340978bb7.png)

[Gunfish slides](http://slides.com/takuyayoshimura-tkyshm/deck-1/fullscreen)

[Gunfish slides (jp)](http://slides.com/takuyayoshimura-tkyshm/deck/fullscreen)

## Quick Started

```bash
$ go get github.com/kayac/Gunfish/cmd/gunfish
$ gunfish -c ./config/gunfish.toml -E production
```

### Commandline Options

option             | required | description
-------------------|----------|------------------------------------------------------------------------------------------------------------------
-port              | Optional | Port number of Gunfish provider server. Default is `8003`.
-environment, -E   | Optional | Default value is `production`.
-conf, -c          | Optional | Please specify this option if you want to change `toml` config file path. (default: `/etc/gunfish/config.toml`.)
-log-level         | Optional | Set the log level as 'warn', 'info', or 'debug'.
-log-format        | Optional | Supports `json` or `ltsv` log formats.
-enable-pprof      | Optional | You can set the flag of pprof debug port open.
-sender-num        | Optional | Set number of concurrency sending notification per http client. That option overwrite config file's param.

## API

### POST /push/apns

To delivery remote notifications via APNS to user's devices.

param | description
--- | ---
Array | Array of JSON dictionary includes 'token' and 'payload' properties

payload param | description
--- | ---
token | Published token from APNS to user's remote device
payload | APNS notification payload

Post JSON example:
```json
[
  {
    "payload": {
      "aps": {
        "alert": "test notification",
        "sound": "default"
      },
      "option1": "foo",
      "option2": "bar"
    },
    "token": "apns device token",
    "header": {
      "apns-id": "your apns id",
      "apns-token": "your app bundle id"
    }
  }
]
```

Response example:
```json
{"result": "ok"}
```

### GET /stats/app

To get the status of APNS proveder server.

stats type | description
--- | ---
Pid | PID
DebugPort | pprof port number
Uptime | uptime of APNS provider server
Workers | number of workers
StartAt | The time of starting APNS provider server
QueueSize | queue size of requests for Gunfish
RetryQueueSize | queue size for resending notification
WorkersQueueSize | summary of worker's queue size
CommandQueueSize | error hook command queue size
RetryCount | summary of retry count
RequestCount | request count to gunfish
ErrCount | count of recieving error response from APNs
SentCount | count of sending notification to APNs

### GET /stats/profile

To get the status of go application.

See detail properties that url: (https://github.com/fukata/golang-stats-api-handler).

## Configuration
The Gunfish configuration file is a TOML file that Gunfish server uses to configure itself.
That configuration file should be located `/etc/gunfish.toml`, and is required to start.
Here is an example configuration:

```toml
[provider]
port = 8203
worker_num = 8
queue_size = 2000
max_request_size = 1000
max_connections = 2000
error_hook = "echo -e 'Hello Gunfish at error hook!'"

[apns]
skip_insecure = true
key_file = "/path/to/server.key"
cert_file = "/path/to/server.crt"

[fcm]
api_key = "API key for FCM"
```

param            | status | description
---------------- | ------ | --------------------------------------------------------------------------------------
port             |optional| Listen port number.
worker_num       |optional| Number of Gunfish owns http clients.
queue_size       |optional| Limit number of posted JSON from the developer application.
max_request_size |optional| Limit size of Posted JSON array.
max_connections  |optional| Max connections
skip_insecure    |optional| Controls whether a client verifies the server's certificate chain and host name.
key_file         |required| The key file path.
cert_file        |required| The cert file path.
error_hook       |optional| Error hook command. This command runs when Gunfish catches an error response.

## Error Hook

Error hook command can get an each error response with JSON format by STDIN.

for example JSON structure: (>= v0.2.x)
```json5
// APNs
{
  "provider": "apns",
  "apns-id": "123e4567-e89b-12d3-a456-42665544000",
  "status": 400,
  "token": "9fe817acbcef8173fb134d8a80123cba243c8376af83db8caf310daab1f23003",
  "reason": "MissingTopic"
}
```

```json5
// FCM
{
  "provider": "fcm",
  "status": 200,
  "registration_id": "8kMSTcfqrca:APA91bEfS-uC1WV374Mg83Lkn43..",
  // or "to": "8kMSTcfqrca:APA91bEfS-uC1WV374Mg83Lkn43..",
  "error": "InvalidRegistration"
}
```

(~ v0.1.x)
```json
{
  "response": {
    "apns-id": "",
    "status": 400
  },
  "response_time": 0.633673848,
  "request": {
    "header": {},
    "token": "9fe817acbcef8173fb134d8a80123cba243c8376af83db8caf310daab1f23003",
    "payload": {
      "aps": {
        "alert": "error alert test",
        "badge": 1,
        "sound": "default"
      },
      "Optional": {
        "option1": "hoge",
        "option2": "hoge"
      }
    },
    "tries": 0
  },
  "error_msg": {
    "reason": "MissingTopic",
    "timestamp": 0
  }
}
```


## Graceful Restart
Gunfish supports graceful restarting based on `Start Server`. So, you should start on `start_server` command if you want graceful to restart.

```bash
### install start_server
$ go get github.com/lestrrat/go-server-starter/cmd/start_server

### Starts Gunfish with start_server
$ start_server --port 38003 --pid-file gunfish.pid -- ./gunfish -c conf/gunfish.toml
```

## Customize

### How to Implement Response Handlers

If you have to handle something on error or on success, you should implement error or success handlers.
For example handlers you should implement is given below:

```go
type CustomYourErrorHandler struct {
    hookCmd string
}

func (ch CustomYourErrorHandler) OnResponse(result Result){
    // ...
}

func (ch CustomYourErrorHandler) HookCmd( ) string {
    return ch.hookCmd
}
```

Then you can use these handlers to set before to start gunfish server `( gunfish.StartServer( Config, Environment ) )`.

```go
InitErrorResponseHandler(CustomYourErrorHandler{hookCmd: "echo 'on error!'"})
```

You can implement a success custom handler in the same way but a hook command is not executed in the success handler in order not to make cpu resource too tight.

### Test
To do test for Gunfish, you have to install [h2o](https://h2o.examp1e.net/). **h2o** is used as APNS mock server. So, if you want to test or optimize parameters for your application, you need to prepare the envronment that h2o APNs Mock server works.

Moreover, you have to build h2o with **mruby-sleep** mrbgem.


```
$ make test
```

### Benchmark
Gunfish repository includes Lua script for the benchmark. You can use wrk command with `err_and_success.lua` script.

```
$ h2o -c conf/h2o/h2o.conf &
$ ./gunfish -c test/gunfish_test.toml -E test
$ wrk2 -t2 -c20 -s bench/scripts/err_and_success.lua -L -R100 http://localhost:38103
```
