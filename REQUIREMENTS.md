# TinkerDB Requirements Spec

## Premise

This is an educational project to build a distributed key-value database from scratch in Go. The project follows a systematic approach to learning distributed systems concepts through hands-on implementation, progressing from a simple standalone server to a fully distributed, production-ready database system.

### Core Objectives
- **Educational Focus**: Learn distributed systems concepts through practical implementation
- **Technology Stack**: Go for backend implementation, with web UI for demonstration
- **Architecture**: Key-value store evolving into a distributed SQL-capable database
- **Timeline**: 12-month development cycle with 6 major milestones


## Milestone 1: The Standalone Core (Months 1-2)

Build a functional key-value server that can store and retrieve data in memory. This will serve as the core for all transactions and basic services in the latter parts of the system.

### Technical Implementation
- **Core Operations**: `SET`, `GET`, `DELETE`, `EXISTS`, `KEYS` (maybe)
- **Protocol**: Decide on a communication protocol
- **Storage**: In-memory hash map with string keys and values, LSM trees, etc.
- **Server**: Go HTTP server with concurrent request handling

### Deliverables
- [ ] Basic key-value server implementation
- [ ] APIs with standard CRUD operations
- [ ] Simple client library (Go)
- [ ] Unit tests for core functionality
- [ ] Basic documentation and usage examples

### Success Criteria
- Server can handle basic CRUD operations
- Multiple clients can connect simultaneously
- Data is stored efficiently in memory
- Implement some form of Tenant separation mechanism.
- Clean, well-documented codebase


## Milestone 2: Durability & Concurrency (Months 3-4)

### Goal
Make the server robust by adding data persistence and proper concurrency control.

### Technical Implementation
- **WAL (Write-Ahead Log)**: Append-only log for all write operations
- **Snapshotting**: Periodic full state dumps for faster recovery
- **Concurrency**: Reader-writer locks for optimal read performance
- **Recovery**: Automatic log replay on server restart

### Deliverables
- [ ] Write-ahead logging implementation
- [ ] Snapshot and recovery mechanisms
- [ ] Thread-safe concurrent access
- [ ] Performance benchmarks (read/write throughput)
- [ ] Crash recovery testing

### Success Criteria
- Data persists across server restarts
- No data corruption under concurrent access
- Fast recovery from crashes
- Maintained performance with persistence

## Milestone 3: Distributed Consensus (Months 5-6)

### Goal
Transform the single server into a fault-tolerant distributed cluster.

### Key Concepts to Learn
- **Raft Consensus Algorithm**: Leader election, log replication, safety
- **State Machine Replication**: Applying operations consistently across nodes
- **Network Programming**: RPC, message serialization, network partitions
- **Fault Tolerance**: Handling node failures and network issues

### Technical Implementation
- **Raft Implementation**: Leader election, log replication, safety guarantees
- **RPC Framework**: gRPC or custom protocol for inter-node communication
- **Cluster Management**: Node discovery, configuration, and health checks
- **Split-Brain Prevention**: Ensuring single leader at all times

### Deliverables
- [ ] Complete Raft consensus implementation
- [ ] Multi-node cluster setup and management
- [ ] Automatic failover and leader election
- [ ] Network partition handling
- [ ] Cluster configuration and deployment scripts

### Success Criteria
- Cluster maintains consistency under normal operation
- Automatic leader election when leader fails
- No split-brain scenarios
- Data consistency across all nodes
- **🎉 First Demo-Ready Version**

## Milestone 4: Automation & Demonstration (Months 7-8)

### Goal
Make the project easily deployable with an interactive demonstration interface.

### Key Concepts to Learn
- **Containerization**: Docker, container orchestration basics
- **Infrastructure as Code**: Terraform, Kubernetes manifests
- **Web Development**: HTML/CSS/JavaScript for demo interface
- **CI/CD**: Automated testing and deployment pipelines

### Technical Implementation
- **Docker**: Multi-stage builds, container optimization
- **Orchestration**: Docker Compose for local development
- **Web UI**: Simple interface for database operations
- **Cloud Deployment**: AWS/GCP/Azure deployment scripts

### Deliverables
- [ ] Docker containerization
- [ ] Docker Compose setup for local development
- [ ] Web-based demo interface
- [ ] Cloud deployment automation
- [ ] CI/CD pipeline setup

### Success Criteria
- One-command deployment to cloud
- Interactive web demo showcasing features
- Automated testing in CI/CD
- Easy onboarding for new contributors


## Milestone 5: Horizontal Scaling (Months 9-10)

### Goal
Enable the database to scale beyond single-machine capacity through sharding.

### Key Concepts to Learn
- **Sharding Strategies**: Range-based, hash-based, directory-based
- **Consistent Hashing**: Minimizing data movement during rebalancing
- **Request Routing**: Smart clients and proxy layers
- **Rebalancing**: Dynamic shard redistribution

### Technical Implementation
- **Sharding Layer**: Consistent hashing for key distribution
- **Router/Proxy**: Request routing to appropriate shards
- **Rebalancing**: Online data migration between shards
- **Cross-Shard Transactions**: Basic distributed transaction support

### Deliverables
- [ ] Sharding implementation with consistent hashing
- [ ] Request routing and load balancing
- [ ] Online rebalancing capabilities
- [ ] Cross-shard query support
- [ ] Shard management tools

### Success Criteria
- Linear scaling with number of shards
- Minimal data movement during rebalancing
- Transparent sharding to clients
- Maintained consistency across shards

## Milestone 6: Production Readiness (Months 11-12)

### Goal
Add production-grade features, monitoring, and comprehensive documentation.

### Key Concepts to Learn
- **Observability**: Metrics, logging, and distributed tracing
- **Performance Optimization**: Profiling, caching, and optimization techniques
- **Security**: Authentication, authorization, and encryption
- **Documentation**: API docs, architecture guides, and contribution guidelines

### Technical Implementation
- **Monitoring**: Prometheus like metrics, structured logging, OpenTelemetry tracing
- **Performance**: Connection pooling, query optimization, caching layers
- **Security**: TLS encryption, basic authentication, RBAC
- **Documentation**: Comprehensive docs, tutorials, and examples

### Deliverables
- [ ] Complete observability stack
- [ ] Performance benchmarks and optimization
- [ ] Security features and best practices
- [ ] Comprehensive documentation
- [ ] Open-source contribution guidelines

### Success Criteria
- Production-ready monitoring and alerting
- Competitive performance benchmarks
- Security best practices implemented
- Welcoming open-source project with clear contribution path

---

## Metrics I want to establish over the course of this project

### Goals
- **Availability**: 99.9% uptime in distributed mode (piggyback off of major cloud provider)
- **Performance**: 10,000+ operations/second per node
- **Latency**: <1ms for local operations, <10ms for distributed operations
- **Consistency**: tbd

---

*This project plan serves as a living document that will evolve as I learn and implement each milestone. The focus remains on education and practical understanding of distributed systems principles.*
