# Consistent Hashing

At this point, we've got the hash ring set up and we're using the consistent hashing technique to place servers and keys
on the ring.

Running `task run` will show that our distribution isn't great. This is because N (number of servers) is quite small.
When we hash the servers and keys, we can easily end up with skew.

To address this, we're going to add virtual nodes. These are references to the actual server node and function to allow
for more even distribution of keys across the hash algorithm's range.

## Steps

1. Update HashRing struct to track how many vnodes to create per server.
1. Update the ring/serverKeys code in AddServer to use vnodes.
1. Update the ring/serverKeys code in RemoveServer to be vnode aware.
