# UCSD CSE224 春季 2025 学期 🔱

[English Version](README.md) | [中文版](README_CN.md)

欢迎来到 UCSD CSE224：**Network Services** 课程仓库！本仓库包含本学期完成的九次实验，涵盖 Go 并发、网络服务和分布式系统开发。请访问各 `lab-*` 目录查看详细实现与说明。


## 🎉 实验 1：Go 入门与项目搭建
- 了解 Go 模块、包管理和基础语法  
- 搭建项目结构，编写第一个 `main.go` CLI 工具  
- 学习 `go build`、`go run` 和 `go test` 命令  

## 🏃‍♂️ 实验 2：并发基础
- 深入学习 goroutine 与 channel  
- 实现生产者-消费者模式，构建并发管道  
- 探索带缓冲与无缓冲 channel 的阻塞与非阻塞行为  

## 🔒 实验 3：同步与竞态
- 使用 `sync.Mutex`、`sync.WaitGroup` 和原子操作保证并发安全  
- 使用 `go run -race` 检测竞态条件  
- 构建并发计数器并保护共享数据结构  

## 🌐 实验 4：限速 FTP 服务
- 基于 `net` 包实现 FTP 服务端与客户端，支持 `STOR` 和 `RETR`  
- 添加下载限速功能，控制传输带宽  

## 🚀 实验 5：gRPC 与 Protocol Buffers
- 使用 `.proto` 文件定义服务，生成 Go 代码  
- 开发 gRPC 客户端与服务器，处理 RPC 请求与响应  
- 理解 HTTP/2 与流控原理  

## 🔄 实验 6：HTTP 服务与 RESTful API
- 使用 `net/http` 包开发 RESTful 接口，处理 JSON 编解码  
- 实现路由、中间件和基础认证  
- 探索 HTTP/2 支持与请求上下文  

## 📊 实验 7：中间件与日志监控
- 构建自定义中间件，实现请求日志、性能监控与限流  
- 集成 `context` 传递请求元数据  
- 监控关键指标：响应时间、吞吐量等  

## 🌎 实验 8：分布式文件存储
- 实现 `RemoveNode` RPC，节点增删时迁移文件并维护一致性  
- 使用一致性哈希分布文件，无需全局锁  
- 确保每个文件仅存在于单个节点  

## 🔑 实验 9：etcd 与集群交互
- 使用 `etcdctl` 与 etcd 集群交互，学习分布式 KV 存储  
- 分析 Leader 选举与节点状态变化  
- 在 AWS 环境中启动/停止节点，观察 IP 变化与选举行为  


欢迎克隆此仓库并深入探索各实验内容。祝编码愉快！🎓  
