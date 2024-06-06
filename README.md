# Fibona-AI-Otel-demo
对用户prompt进行链路追踪，监控prompt到LLM的过程  

## 运行所需条件
* proto编译器: 用于编译proto文件
* opentelemetry-collector: 用于数据的收集和导出 运行时带上otel-config.yaml配置文件
* jaeger： 导入数据，进行可视化
## 所用技术栈
### LLM
* 采用langchain本地部署大模型
* 使用Faiss构建本地知识库
### golang-client
* gin: 实现web服务
* gorm: 实现数据的存储
* proto： 编写proto文件，用protoc命令生成相应的go文件，用于grpc服务
* grpc: 实现python的远程方法调用
* otel: 使用了otel golang sdk 实现对数据的链路追踪和监控，并用otel收集器将span收集并导出到jaeger进行可视化分析
### python-server
* proto: 和golang中的proto文件保持一致，用protoc命令生成相应的go文件，用于grpc服务
* grpc： 实现远程调用的方法，并将方法注册到grpc中，供golang-client调用
* otel： 与golang中相似，将调用本地大模型过程和从向量数据库中提取信息的过程进行监控并导出到jaeger
* langchain： 使用与langchain生态的各种库，例如sentence_transformers，RecursiveCharacterTextSplitter，FAISS 等等，将向量数据库中的数据加载出来并进行分词处理，选取符合的知识发送给LLM，使输出结果更加符合用户propmt。 

## jaeger监控效果图
![alt text](/assert/image-1.png)
### golang-client监控
父sapn：调用python远程服务
![alt text](/assert/image-2.png)
子span： 将数据存储在mysql中
![alt text](/assert/image-3.png)
###  python-server监控
父span：将用户propmt发送给LLM过程
![alt text](/assert/image-4.png)
子span： 向量数据库检索过程
![alt text](/assert/image-5.png)
## MYSQL中的数据
![alt text](/assert/image-6.png)
