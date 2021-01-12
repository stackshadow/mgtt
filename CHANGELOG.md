## v0.20.0 - AuthAPI

### New
- [auth] Add MQTT-Topic: '$SYS/auth/user/+/set'

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
