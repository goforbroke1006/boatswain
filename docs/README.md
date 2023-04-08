# Boatswain

### Data flow by components

Node:

- <- "boatswain/transaction" topic
- <- "boatswain/consensus-vote" topic
- -> "boatswain/consensus-vote" topic
- -> "boatswain/reconciliation/req" topic
- <- "boatswain/reconciliation/resp" topic

DApp Chat:

- <- "chat: <CHAT ID>" topic
- -> "chat: <CHAT ID>" topic
- -> "boatswain/transaction" topic
- <- "boatswain/reconciliation/resp" topic

### Share new data inside topic "chat/cat-owners-worldwide"

```json
{
  "blockchain": "",
  "id": "<uuid>",
  "peer_sender": "peer-1",
  "peer_recipient": "",
  "content": "Hi there!",
  "timestamp": 1680950512
}
```

### Put transaction into node with "boatswain/transaction" topic

```json
{
  "blockchain": "chat/cat-owners-worldwide",
  "id": "uuid",
  "peer_sender": "peer-1",
  "peer_recipient": "",
  "content": "Hi there!",
  "timestamp": 1680950512
}
```

### Vote for consensus with "boatswain/consensus-vote" topic

```json
{
  "blockchain": "chat/cat-owners-worldwide",
  "index": 123456,
  "hash": "bbbbbbbbbbb",
  "previousHash": "aaaaaaaaaaa",
  "timestamp": 1680950612,
  "data": [
    {
      "id": "05d791bc-1846-407b-ae55-070266dad475",
      "peer_sender": "peer-1",
      "peer_recipient": "",
      "content": "Hi there!",
      "timestamp": 1680950512
    },
    {
      "id": "b2350be2-8270-494d-9eeb-2a5851c2fbed",
      "peer_sender": "peer-1",
      "peer_recipient": "",
      "content": "Where I can buy turkish angora in Europe?",
      "timestamp": 1680950523
    }
  ]
}
```

### Sync blocks on node start flow

```mermaid
sequenceDiagram
    participant BlockStorage
    participant Syncer
    participant Reconciliation Req
    participant Reconciliation Resp
    participant Another Node
    participant Another BlockStorage

    loop If BlockStorage empty
        Syncer ->> BlockStorage: create genesis
    end

    loop Till Reconciliation Resp returns non-empty payload
        Syncer ->> BlockStorage: load last block
        BlockStorage -->> Syncer: return last block

        Syncer ->> Reconciliation Req: request N blocks after local last block

        Reconciliation Req ->> Another Node: catch request
        Another Node ->> Another BlockStorage: load N block after block X
        Another BlockStorage -->> Another Node: return N block after block X
        Another Node ->> Reconciliation Resp: write N blocks

        Reconciliation Resp ->> Syncer: return next N blocks after local last block
        Syncer ->> BlockStorage: append blocks
    end
```

### Normal data exchange flow

```mermaid
sequenceDiagram
    participant Blockchain 1
    participant Node 1
    participant Chat Peer 1
    participant Chat Peer N
    participant Node N
    participant Blockchain N

    Chat Peer 1 ->> Node 1: transaction "Hi there!"
    Chat Peer 1 ->> Chat Peer N: transaction "Hi there!"
    Chat Peer N ->> Node N: transaction "Hello!"
    Chat Peer N ->> Chat Peer 1: transaction "Hello!"

    loop Node 1 collect message for 1 Mb check
        Node 1 ->> Node N: Vote for consensus
    end

    loop Node N collect message for 1 Mb check
        Node N ->> Node 1: Vote for consensus
    end

    loop Node 1 vote timeout
        Node 1 ->> Blockchain 1: store block
    end

    loop Node N vote timeout
        Node N ->> Blockchain N: store block
    end
```