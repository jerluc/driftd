# Under development

* **Create Unix socket interface** that would allow external software to
  interact with riftd without needing to have any special libraries
  installed. Something as simple as a domain socket or a managed
  directory of sockets that would allow riftd to forward
  incoming/outgoing messages to/from various remote endpoints is all
  that's necessary.
* **Implement simple routing/addressing scheme** that would allow remote
  nodes to be distinguishable. For now, this address will simply be the
  64-bit device MAC address, but maybe something more logical would be
  appropriate? Remember that the addressing scheme does not necessarily
  need to be transitive, and could very easily be a temporal thing.
* **Implement ARP-style broadcast/unicast protocol** to allow for node
  discovery within range.
* **Research the creation of a simple Linux kernel module** that would
  allow us to implement the protocol as a proper socket protocol. This
  would then allow any program that can use the `socket()` syscall to
  create a new Rift socket.
