# Fixing the Naive Hashing

We've got a working (modulo-based) hash in place, but it's got issues. We've currently got a map and a slice holding
servers that have been added to the system.

In this lab, we're going to adapt it to use a hash ring (map of uint32 server hashes to server), and a sorted list of
server hashes.

## Steps

Check out _hashring/hashing.go_. You'll see some TODO comments in there. Specifically, we're going to:

1. Update HashRing struct to track `ring`, `servers`, and `serverKeys`.
1. Update AddServer to handle insertion.
1. Update RemoveServer to handle removal.
1. Update GetServer to use the new fields.
1. Anything else required to make the tests pass again.
