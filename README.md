container-updater
=================
Automatically Executes commands when a request arrives on the http interface. The container comes with a docker client preinstalled as to make it easier to control the host instance, this requires a host socket to be mounted, e.g:

Based on the idea of [autodock](dockerun -v /var/run/docker.sock:/var/run/docker.sock -p 30000:30000 cu 'docker version') but aims to improve it in several ways:

- Simplify the webhook receiving part by assuming we can just spin up a new container for each webhook we want to capture.
- Added a runner configuration that allows the configuration of the shell environment, defaults to `sh -c {{cmd}}`
- Improved argument-to-command parsing and made it easier to specify several commands after each other
- added unit tests for argument parsing

```
docker run -v $HOSTPATH_TO_SOCKET:/var/run/docker.sock -p 30000:30000 cu 'docker pull'
```

On a typical Mac OSX boot2docker instance this might look like:
```
docker run -v /var/run/docker.sock:/var/run/docker.sock \
-p 30000:30000 cu \
'docker pull ubuntu' \
'docker run -p 8080:8080 ubuntu:14:04 echo "hello world"'
```
This will, whenever a http request arrives at port 30000, (1) pull the latest ubuntu images from the default repository and (2) when done spin up a new instance that exposes port 8080 and echos "hello world"
