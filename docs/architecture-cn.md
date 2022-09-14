# Architecture

![comparison diagram](./Architecture-m1.png)

FileDag Storage Milestone 1的整体架构是分层的，由三个部分组成:存储层、池层和对象层。每一层实现都是独立的、灵活的、可持续迭代的，因为它们之间只通过RPC协议进行交互。

我们来看看最底层Storage layer，它主要关注数据块的存储，解决快速高效的读写问题。目前采用的是键值存储方案。

中间层Pool层主要管理DAG结构化数据的读写和访问权限控制。它管理DagNode集群。每个DagNode管理多个datanode，采用Erasure Coding分片技术。每个片段通过DataNode Client存储到相应的DataNode服务中。已经根据Milestone 1实现了单个DagNode方案。DagPool客户端实现了块存储接口，可以在任何需要块存储的地方使用。

最后，最上层对象层实现对象存储的全部功能，主要分为S3、IAM和Store。因此，FileDag Storage兼容S3接口，具有身份权限控制。
