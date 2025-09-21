# Initial implementation

We're going to create a basic modulo-based hash to handle balancing keys across servers using the modulo hash pattern.
The pitfalls of this approach were outlined earlier, but IMO it's worth building up from this to a fully-fledges hash
ring with virtual nodes.

> Running `task run` right now will effectively be useless, since there's no keys added (not implemented yet), so the
> distribution isn't measurable.

## Steps

For each step here, be sure to add/update tests so we're able to safely refactor later. While my preference would be
test first, there's nothing complex here that would make writing them after harder or less effective. Dealer's choice!

1. Together, we'll implement `NewHashRing` and the `hashKey` function.
1. Implement `Size` and `AddServer` methods (return error if server already exists).
1. Implement `RemoveServer` method (return error if server not found).
1. Implement `GetServer` and `GetServers` methods.
1. Implement `GetDistribution`, and flesh out `TestConsistency`.

Running `task run` now, should show a working demo. Notice the distribution and the number of keys that need to move
when a server is added/removed. This will be the focus of the rest of the lab.
