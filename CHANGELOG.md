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
