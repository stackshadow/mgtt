## v0.34.0

### Fixes

- 🐛 [broker] Fix receive of own package

## v0.33.0

### New

- [broker] New flag `--enable-admin-topics` to enable admin-topics
- [broker] Add `$SYS/ping`-Topic that just respond an `$SYS/pong`
- [auth] Topic `$SYS/auth/user//get` ( no username given ) return now the infos for `_anonym`-User

### Breaking changes

- no breaking changes

### Fixes

- [auth] Prevent set of password for `_anonym`-User

### Under the hood

- [auth] Fix tests
- [clientlist] Use auto-init of package
- [cli] Use auto-init of package
- [ci] Release-jop only created on tags
- [logging] Print out debug-level
- [logging] Show caller on every log-level

## v0.32.0

### New

- [broker] $SYS-Topics will not be published to all clients ( only to sending-client )
- [broker] Serve() will now return when server is listening

### Breaking changes

- [broker] $SYS/broker/version renamed to $METRIC/broker/version

## v0.31.0 - Read the manual

- [doc] Update the whole documentation
- Update to go 1.16

## v0.30.0 - Session affinity

### New

- [broker] Session functionality
- [broker] Add unsubscribe, Closes #3
- [auth] $SYS/auth/user/%s/set/success return now the user-object
- [auth] userSet return now the modified user-object

### Fixes

- [auth] $SYS/auth/users/list/get remove userless check of existing broker
- [auth] Add configUserGet to get the userinfo in a secured way
- [auth] Write default config-file if no file is present
- [auth] Add groups-block to default config-file
- [clientlist] Make clientlist thread-save
- [cli] Security fixes for files

### under the hood

- [broker] Add test for client-package
- [broker] Rename handlePacket -> onPacket
- [plugins] Add DeRegister() to unregister an plugin
- [tests] Move mocked-network-connection to dedicated package
- [client] Move machter-function to separate file
- [build] Add gosec-badge
- [ci] Add version label to docker-image
- [cli] Add command to create list of environment variables for documentation

## v0.20.0 - AuthAPI

### New

- [auth] Add MQTT-Topic: '$SYS/auth/user/+/set'
- [auth] $SYS/auth/user/+/get includes now the username in the user-object

### Fixes

- [auth] Dont send password in the list of users
- [auth] Json settings for user-object
- [auth] Fix tests and adjust according to breaking changes
- [auth] Update documentation according to the new topix

### Breaking changes

- [auth] Change structure of config-file
- [auth] Rename '$SYS/auth/user/%s/delete/ok' to '$SYS/auth/user/%s/delete/success'
- [auth] Rename MQTT-Topic: '$SYS/self/username' to '$SYS/self/username/string'
- [auth] Remove MQTT-Topic: '$SYS/auth/user/+/password/set' in favor for '$SYS/auth/user/+/set'
- [auth] Remove MQTT-Topic: '$SYS/self/username/get' and '$SYS/self/groups/get' in favor for '$SYS/self/user/get'

### CI

- Add: coverage-badge
- Add automatic tests and coverage-badge
- run tests before build of docker-image
- Add missing gawk

## v0.16.0

- [gitlab-ci] Add gocyclo-badge
- [gitlab-ci] Add lastbuild-badge
- [cli] Fix cli-default environment-parameter
- [cli] Add env-var: SELFSIGNED environment-var
- [auth] Add env-var: AUTH_USERNAME and AUTH_PASSWORD. this create a new user with the specified password
- [auth] Add env-var: AUTH_ANONYMOUSE enable anonymouse auth
- [auth] Add: get of an user '$SYS/auth/user/+/get'
- [auth] Add: get of the user-group '$SYS/self/groups/get'
- [plugins] Add OnDisconnected()-callback to inform plugins about an client-disconnect

### Broken changes

- [auth] Rename: '$SYS/auth/users/list' to '$SYS/auth/users/list/get'

## v0.14.0

- [clientlist] Add new package clientlist that holds all our connected clients
- [broker] Add Stop of server
- [broker] Rework broker-loops
- [client] Remove receiver-loop
- [client] Add sender-loop
- Fix all tests

## v0.13.0

### REFACTOR

- simplify message store for packets
- save pubrecs to memory not to db
- simplyfy QoS2
- add tests for QoS1 and QoS2
