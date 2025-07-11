# UCSD🔱 CSE224 Spring 2025 Quarter

Welcome to the UCSD CSE224: **Network Services** course repository for Spring 2025! This repository contains all nine labs completed during the course, covering Go concurrency, network services, and distributed systems development. Explore each `lab-*` directory for detailed implementations and instructions.


## 🎉 Lab 1: Go Introduction & Project Setup
- Introduction to Go modules, package management, and basic syntax  
- Scaffolded the project structure and wrote the first `main.go` CLI tool  
- Learned `go build`, `go run`, and `go test` commands  

## 🏃‍♂️ Lab 2: Concurrency Basics
- Deep dive into goroutines and channels  
- Implemented the producer-consumer pattern and built concurrent pipelines  
- Explored buffered vs. unbuffered channels and blocking behavior  

## 🔒 Lab 3: Synchronization & Race Conditions
- Used `sync.Mutex`, `sync.WaitGroup`, and atomic operations for safe concurrency  
- Detected race conditions with `go run -race`  
- Built a concurrent counter and protected shared data structures  

## 🌐 Lab 4: Rate-Limited FTP Service
- Built an FTP server and client using the `net` package, supporting `STOR` and `RETR`  
- Added optional download rate limiting to control bandwidth  

## 🚀 Lab 5: gRPC & Protocol Buffers
- Developed gRPC clients and servers in Go, defining services in `.proto` files  
- Generated Go code from Protocol Buffers and implemented RPC handlers  
- Learned about HTTP/2 and flow control in gRPC  

## 🔄 Lab 6: HTTP Server & RESTful API
- Created a RESTful API with the `net/http` package, handling JSON encoding/decoding  
- Implemented routing, middleware, and basic authentication  
- Explored HTTP/2 support and request context usage  

## 📊 Lab 7: Middleware & Logging
- Built custom middleware for request logging, performance monitoring, and rate limiting  
- Integrated `context` for propagating metadata  
- Monitored key metrics such as response time and throughput  

## 🌎 Lab 8: Distributed File Storage
- Implemented `RemoveNode` RPC to migrate files when nodes join or leave  
- Used consistent hashing to distribute files across nodes without global locks  
- Ensured file consistency and single-copy guarantees per lab specifications  

## 🔑 Lab 9: etcd & Cluster Interaction
- Interacted with an etcd cluster using `etcdctl`, learning distributed key-value storage  
- Analyzed leader election and node state transitions  
- Deployed nodes on AWS, observed IP changes, and studied election behaviors  


Feel free to clone this repository and dive into each lab for more details. Enjoy! 🎓
