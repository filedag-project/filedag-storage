# FileDAG Storage

<!-- ABOUT THE PROJECT -->
## About The Project
FileDAG Storage A distributed storage service built on the IPFS technology stack. Different from the official implementation of IPFS, we focus more on data management, data reliability, availability and fault tolerance, and clustering of storage nodes.
 
The minimum storage unit of FileDAG Storage is dag, which is the data block. Files or objects are organized in a merkle-dag structure, and multiple files or objects may share some data blocks. The obvious benefit of this is to reduce redundant data, especially for multi-version systems. Not only does reduce data redundancy, it also saves bandwidth on network transmissions.  Every advantage has its disadvantage, the downside is that data management has become more complex. First, the file or object management module needs to be abstracted on the basis of merkle-dag; second, the file cannot be deleted directly, only the data blocks that are no longer needed can be released through garbage collection. The situation becomes more complicated when storage resources are used in a multi-user fashion. The current implementation of IPFS does not support multiple users.

Supporting multiple users and maintaining a larger dag pool will magnify the advantages of reducing redundant data and saving network bandwidth, while also increasing the difficulty of data management. For development teams that use the IPFS technology stack to build storage services and participate in the Filecoin distributed storage network, this is an unavoidable challenge.

To provide commercial storage services, we must first ensure data reliability, availability and fault tolerance. IPFS is currently implemented based on a single node, it is relatively lightweight and suitable for personal use.

The development of FileDAG Storage will provide a solution to the above challenges
- Learn from object storage services about how to do data management
- Handle the release of data blocks in a reference counting manner
- Provide data availability by using distributed storage nodes, and use erasure coding technology to improve data reliability and fault tolerance
  


## Architecture

- Dag Node - The underlying storage module, responsible for the storage and release of data blocks, similar to blockstore
- Dag Pool - A virtual data block set composed of distributed Dag Nodes, supports multi-user use, and provides authentication methods；Responsible for the implementation of clustered storage nodes and provide data fault tolerance solutions；Responsible for the reference record of data blocks, providing an interruptible garbage collection mechanism；
- Object Store - Object storage service abstraction layer based on Dag Pool, responsible for partially providing s3-compatible api interfaces;
- Control Pannel - Visual management ui


## Roadmap

- [ ] The development of a single Dag Node, which satisfies the interface of blockstore, provides basic storage services for Dag Pool
- [ ] Implement the multi-user access and authentication mechanism of Dag Pool
- [ ] Implement Dag Pool's reference record to data blocks
- [ ] Implement Dag Pool's data pin strategy
- [ ] Implement the interruptible garbage collection mechanism of Dag Pool


- [ ] Dag Node - Develop data fault tolerance based on Reed-Solomon Erasure Code
- [ ] Dag Pool - Organize multiple Dag Nodes to build a storage cluster based on libp2p and redis hash slots
- [ ] Dag Pool - Storage node health report and global consistency state implementation
- [ ] Dag Pool - Dynamic expansion of storage nodes
- [ ] Dag Pool - Dynamic scaling of storage nodes



- [ ] Object Store - Implementation of basic data structures such as users, regions, buckets, and objects
- [ ] Object Store - Implement user authentication api
- [ ] Object Store - Implement bucket related operations api
- [ ] Object Store - Implement object manipulation api 
- [ ] Object Store - Implement permission operation api


- [ ] Control Pannel - Implement Dashboard - storage pool statistics overview 
- [ ] Control Pannel - Implement user interface
- [ ] Control Pannel - Implement Object Store user and access interface
- [ ] Control Pannel - Implement the object store bucket operation interface
- [ ] Control Pannel - Implement Object Store permission setting interface


- [ ] Implement satellite nodes that can connect to the IPFS network in the outer layer of Dag Pool
- [ ] Provide lightweight IPFS gateway services according to user customization



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


