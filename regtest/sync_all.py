#!/usr/bin/python3
from bitcoinrpc.authproxy import AuthServiceProxy, JSONRPCException
import time

'''
addnode=127.0.0.1:11156 #0
addnode=127.0.0.1:11157 #1
addnode=127.0.0.1:11158 #2
addnode=127.0.0.1:11159 #3
addnode=127.0.0.1:11160 #4
addnode=127.0.0.1:11161 #5
addnode=127.0.0.1:11162 #6
addnode=127.0.0.1:11163 #7
addnode=127.0.0.1:11164 #8
addnode=127.0.0.1:11165 #9
addnode=127.0.0.1:11166 #10
addnode=127.0.0.1:11167 #11
addnode=127.0.0.1:11168 #12
addnode=127.0.0.1:11169 #13
#addnode=127.0.0.1:11170 #14
addnode=127.0.0.1:11171 #15
addnode=127.0.0.1:11172 #16
addnode=127.0.0.1:11173 #17
addnode=127.0.0.1:11174 #18
'''
ports = {
    0 : 12156,
    1 : 12157,
    2: 12158,
    3: 12159,
    4: 12160,
    5: 12161,
    6: 12162,
    7: 12163,
    8: 12164,
    9: 12165,
    10: 12166,
    11: 12167,
    12: 12168,
    13: 12169,
    14: 12170,
    15: 12171,
}

keys = ports.keys()
minerId = 13
nodes = []

# generate nodes
for id in keys:
    nodes.append(AuthServiceProxy("http://%s:%s@127.0.0.1:%d"%("rt", "rt", ports[id])))

def test():
    for id in keys:
        best_block_hash = nodes[id].getbestblockhash()
        print(nodes[id].getblock(best_block_hash))


def sync_blocks(rpc_connections, wait=1, stop_after=-1):
    """
    Wait until everybody has the same block count
    """
    print("Waiting for blocks to sync (wait interval=%d sec each, max tries=%d" %(wait, stop_after))
    count = 0
    while True:
        counts = [ x.getblockcount() for x in rpc_connections ]
        count += 1
        if counts == [ counts[0] ]*len(counts):
            break
        if stop_after != -1 and count > stop_after:
            break
        print("loop = " + str(count))
        time.sleep(wait)

def sync_mempools(rpc_connections, wait=1, stop_after=-1):
    """
    Wait until everybody has the same transactions in their memory
    pools
    """
    print("Waiting for mempools to sync (wait interval=%d sec each, max tries=%d" %(wait, stop_after))
    count = 0
    while True:
        pool = set(rpc_connections[0].getrawmempool())
        num_match = 1
        count += 1
        for i in range(1, len(rpc_connections)):
            if set(rpc_connections[i].getrawmempool()) == pool:
                num_match = num_match+1
        if num_match == len(rpc_connections):
            break
        if stop_after != -1 and count > stop_after:
            break
        print("loop = " + str(count))
        time.sleep(wait)

test()
sync_blocks(nodes)
sync_mempools(nodes)

nodes[minerId].generate(1)

sync_blocks(nodes)
sync_mempools(nodes)
