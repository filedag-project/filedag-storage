# FileDAG Storage

<!-- ABOUT THE PROJECT -->
## About The Project
FileDAG Storage 是基于 IPFS 技术栈来构建的分布式存储服务。区别于IPFS的官方实现，我们更关注于数据的管理，数据的可靠性、可用性和容错性，以及存储节点集群化。

FileDAG Storage 的最小存储单位是 dag, 即为数据块 block。文件或对象以 merkle-dag 的结构来组织，多个文件或对象可能会共享部分数据块。这样带来的明显好处是减少了冗余数据，对于多版本系统的去冗余效果尤其明显。不仅仅是存储减少了冗余，在网络传输上同样节省了带宽。凡事有利必有弊，不好的地方在于数据管理变得更加复杂了。一则，需要在 merkle-dag 的基础上抽象出文件或对象管理模块；二则，无法直接删除文件，只有通过垃圾回收的方式来释放不再需要的数据块。如果以多用户方式来使用存储资源，情况会变得更加复杂。IPFS 目前的实现是不支持多用户的。

支持多用户，维护更大的 dag 存储池，会放大减少冗余数据和节省网络带宽的优势，同时也放大了数据管理的难度。对于利用 IPFS 技术栈来组建存储服务，参与到 Filecoin 分布式存储网络中的开发团队来讲，这是一件难以回避的事情。

提供商业化的存储服务，首先要保证数据的可靠性、可用性和容错性。IPFS 目前是基于单节点来实现的，没有考虑这些，比较轻量，适合个人使用。

FileDAG Storage 的开发将为上述问题提供提供一种解决方案。
- 学习成熟的对象存储服务的技术方案来管理数据
- 以引用计数的方式来处理数据块的释放
- 通过使用分布式的存储节点来提供数据的可用性，利用纠删码的技术来提高数据的可靠性和容错性
  


## Architecture

- Dag Node - 底层的存储模块, 负责数据块的存储和释放，类似于 blockstore
- Dag Pool - 由分布式的 Dag Node 构成的一个虚拟的数据块集合, 支持多用户使用，提供认证方式；负责存储节点集群化的实施，提供数据容错方案；负责数据块的引用记录，提供一种可中断的垃圾回收机制；
- Object Store - 基于 Dag Pool 组建的对象存储服务抽象层，负责提供部分兼容 s3 的 api 接口;
- Control Pannel - 可视化管理界面


## Roadmap







<!-- CONTRIBUTING -->
## Contributing

PRs are welcome!



<!-- LICENSE -->
## License

Distributed under the MIT License. 



<!-- CONTACT -->
## Contact




<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements


