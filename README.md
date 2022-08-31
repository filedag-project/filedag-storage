# FileDAG Storage

[![LICENSE](https://img.shields.io/github/license/filedag-project/filedag-storage)](./LICENSE "LICENSE")
[![Build Status](https://img.shields.io/github/workflow/status/filedag-project/filedag-storage/Go)]()

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

### Milestone 1

#### Goal:

Build fundamental data structure, and the overall architecture of this project.

#### Description:

- Development of single DAG Node:
    - [x] supports API of the block store
    - [x] providers basic storage service for DAG Pool
- DAG Pool:
    - [x] multi-user access
    - [x] authentication mechanism
- Object store:
    - [x] implements basic data structure, such as user, region, bucket, and object
    - [x] implements API of user authentication

### Millstone 2

#### Goal: 

Implement data management.

#### Description:

- DAG Pool:
    - reference records of data blocks
    - strategy of data pin
    - interruptible garbage collection mechanism of DAG Pool
- Object store:
    - API of bucket related operations
    - API of object manipulation
    - API of permission operation


### Milestone 3

#### Goal: 

Realize clustered DAG Node and development of data fault tolerance.

#### Description:

- DAG Node:
    - develops data fault tolerance based on Reed-Solomon Erasure Code
- DAG Pool:
    - organizes multiple DAG Nodes to build a storage cluster based on libp2p and Redis Hash Slots
    - provides health report of storage nodes and status of global consistency
    - supports dynamic expansion of storage nodes
    - supports dynamic scaling of storage nodes

### Milestone 4

#### Goal: 

Develop the Control Panel of FileDAG Storage.

#### Description:

- Dashboard of storage pool statistics overview
- User interface
- Object Store interface:
    - users
    - user access operations
    - bucket operations
    - permission setting
    
### Milestone 5

#### Goal: 

Connect with IPFS.

#### Description:

- Implements Satellite Nodes connected to the IPFS network in the outer layer of the DAG Pool
- Provides lightweight IPFS gateway services according to user customization


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


