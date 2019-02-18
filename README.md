# GoWinter

Solution for [Talk to Zombies](https://github.com/mysteriumnetwork/winter-is-coming) challenge 
and learning Golang on the way.

Tested against `go version go1.11.4 darwin/amd64`

# Running tests

`$ go test ./...`

# Running the game server

```
$ go build
$ ./gowinter
```

# Connecting as a client

```
$ telnet localhost 52000
```

# Client commands

| Command               | Description                        |
| --------------------- | ---------------------------------- |
| `START <player-name>` | Starts the game with `player-name` |
| `SHOOT <x> <y>`       | Shoots at `(x, y)` | 
