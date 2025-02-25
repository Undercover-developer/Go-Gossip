## How to Run the Nodes
**Spin up the Bootstrap Node**:
Start the first node (which will act as the bootstrap node) without specifying a bootstrap address. For example:

```bash
go run main.go -node=1
```
This node will listen on port 8000 + 1 = 8001 and will maintain its own peer list.

**Spin up Peer Nodes**:
For each additional node, pass the bootstrap node’s address so that they can join the network. For example, to start two more nodes:

```bash
go run main.go -node=2 -bootstrap=127.0.0.1:8001
```
```bash
go run main.go -node=3 -bootstrap=127.0.0.1:8001
```
These nodes will have addresses 127.0.0.1:8002 and 127.0.0.1:8003, respectively. They will contact the bootstrap node to join the network.

**Gossiping**:
Once joined, each node will periodically select a random subset of peers (based on the provided fanout) and send gossip messages.