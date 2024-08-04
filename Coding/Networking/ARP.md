# Address Resolution Protocol
Used to discover [[MAC Address]]es of devices on the network and map them to an associated [[IP Address]].

Gateways such as switches or routers are [[Layer 2]] devices which require a [[MAC Address]] for communication as opposed to a [[Layer 1]] device which only requires an [[IP Address]].

EX: Host A needs to find the [[MAC Address]] of the gateway (usually 10.0.0.1 or 192.168.1.1)

1. Host A [[Broadcast]]s a [[Frame]] to the whole network asking "Who has [[MAC Address]] of *[[IP Address]]*, tell *sending [[IP Address]]*".

2. Host B sees the message but realizes it's not the correct recipient and discards it.

3. The gateway sees the message and replies with it's [[MAC Address]].

4. Host A stores the received message in it's [[ARP Cache]] associating it with an [[IP Address]].

5. Host A wants to send a request to an external server, it accesses the gateway's location through the [[ARP Cache]] and sends the request to the gateway for handling.

# [[ARP]] [[Frame]]s

## Ethernet Frame Header

| # of bits | name                        | ex                |
| --------- | --------------------------- | ----------------- |
| 48        | destination [[MAC Address]] | ff:ff:ff:ff:ff:ff |
| 48        | source [[MAC Address]]      | 00:00:00:00:00:00 |
| 16        | ether type ([[ARP]])        | 0x0806            |
|           |                             |                   |

## [[ARP]] Payload

| # of bits | name                      | ex                                         |
| --------- | ------------------------- | ------------------------------------------ |
| 16        | hardware type             | 0x0001 (Ethernet)                          |
| 16        | protocol type             | 0x0800 (ipv4)                              |
| 8         | hardware size             | 0x06                                       |
| 8         | protocol size             | 0x04                                       |
| 16        | opcode                    | 0x0001 (ARP Request) or 0x0002 (ARP Reply) |
| 48        | [[MAC Address]] of sender | 70:5E:DB:27:24:3D                          |
| 32        | [[IP Address]] of sender  | 192.168.0.18                               |
| 48        | [[MAC Address]] of target | 00:00:00:00:00:00 (empty)                  |
| 32        | [[IP Address]] of target  | 192.168.0.19                               |
|           |                           |                                            |

Ethernet Frame Header + [[ARP]] Payload = complete & valid [[ARP]] [[Frame]]

