### WFS文件存储系统，主要解决海量小文件的存储问题 [[English document]](https://github.com/donnie4w/wfs/blob/main/README.md "[English document]")

###### wfs有非常高效的读写效率，在高并发压力下，wfs存储引擎读写响应时间可以达到微秒级别.

##### 海量小文件可能带来的许多问题：

   海量小文件存储在不同的硬件环境和系统架构下，都会带来一系列显著的问题。无论是采用传统的机械硬盘（HDD）还是现代的固态硬盘（SSD），这些问题都可能影响系统的性能、效率、扩展性和成本：

1. 存储效率低下：对于任何类型的硬盘，小文件通常会导致物理存储空间的低效使用。由于硬盘有其最小存储单元（扇区或页），小文件可能会占用超过其实际内容大小的空间，尤其是在每个文件还需额外存储元数据的情况下，如inode（在Unix-like系统中）或其他形式的元数据记录，这会进一步加大空间浪费。inode耗尽：每个文件和目录至少占用一个inode，而inode的数量是在格式化磁盘并创建文件系统时预先设定的。当系统中有大量小文件时，即使硬盘空间还很充足，也可能因为inode用完而导致无法继续创建新文件，尽管剩余磁盘空间足以存放更多数据。性能影响：随着inode数量增多，查找和管理这些inode所对应的元数据会变得更复杂和耗时，尤其是对于不支持高效索引机制的传统文件系统，这会影响文件系统的整体性能。扩展性受限：文件系统设计时通常有一个固定的inode总数，除非通过特殊手段（如调整文件系统或重新格式化时指定更多inode），否则无法动态增加inode数量来适应小文件增长的需求。
2. I/O性能瓶颈与资源消耗：在HDD环境中，随机读写大量小文件会引发频繁的磁盘寻道操作，从而降低整体I/O性能。而在SSD中，尽管寻道时间几乎可以忽略，但过于密集的小文件访问仍可能导致控制器压力增大、写入放大效应以及垃圾回收机制负担加重。
3. 索引与查询效率问题：海量小文件对文件系统的索引结构形成挑战，随着文件数量的增长，查找、更新和删除小文件时所需的元数据操作会变得非常耗时。尤其在需要快速检索和分析场景下，传统索引方法难以提供高效的查询服务。
4. 备份恢复复杂性与效率：备份海量小文件是一个繁琐且耗时的过程，同时在恢复过程中，尤其是按需恢复单个文件时，需要从大量备份数据中定位目标文件，这将极大地影响恢复速度和效率。
5. 扩展性与可用性挑战：存储系统在处理海量小文件时，可能面临扩展性难题。随着文件数量的增长，如何有效分配和管理资源以维持良好的性能和稳定性是一大考验。在分布式存储系统中，还可能出现热点问题，导致部分节点负载过高，影响整个系统的稳定性和可用性。

** wfs 作用在于将海量提交存储的小文件进行高效的压缩归档。并提供简洁的数据获取方式，以及后台文件管理，文件碎片整理等。**

------------

#### wfs相关程序

- wfs源码地址       https://github.com/donnie4w/wfs
- go客户端           https://github.com/donnie4w/wfs-goclient
- java客户端         https://github.com/donnie4w/wfs-jclient
- python客户端    https://github.com/donnie4w/wfs-pyclient
- wfs在线体验      http://testwfs.tlnet.top     用户名 admin     密码 123
- wfs使用文档      https://tlnet.top/wfsdoc

------------

#### wfs的特点

- 高效性
- 简易性
- 零依赖
- 管理平台
- 图片处理

------------

#### 应用场景

- 媒体存储：适用于存储和访问海量的小文件，如图片、文本等。凭借高性能存储引擎，WFS 可实现高速存取，并提供丰富的图片资源处理功能。

------------

#### 技术特点

- 高吞吐量低延迟：保证在高并发场景下的数据存取速度。
- 支持多级别数据压缩存储：节省存储空间，提高存储效率。
- 支持http(https)协议存取文件
- 支持thrift协议长连接存取文件
- 支持图片基本处理：内置图片处理功能，满足多媒体存储需求。

------------

#### WFS的压力测试与性能评估

###### 请注意，以下基准测试数据主要针对WFS数据存储引擎，未考虑网络因素的影响。在理想条件下，基于基准测试数据得出估算数据

**以下为部分压测数据截图**
![](https://tlnet.top/f/1709371893_7752.jpg)

![](https://tlnet.top/f/1709371933_7249.jpg)

![](https://tlnet.top/f/1709373380_17625.jpg)

![](https://tlnet.top/f/1709373414_15548.jpg)

##### 测试数据说明：

- 第一列为测试方法，写Append, 读Get ， *-4四核，*-8八核
- 第二列为本轮测试执行总次数
- ns/op: 每执行一次消耗的时间
- B/op：每执行一次消耗的内存
- allocs/op：每执行一次分配内存次数

##### 根据基准测试数据，可以估算出wfs存储引擎的性能

- 存储数据性能估算
1. Benchmark_Append-4 平均每秒执行的操作次数约为：1 / (36489 ns/operation) ≈ 27405次/s
2. Benchmark_Append-8 平均每秒执行的操作次数约为：1 / (31303 ns/operation) ≈ 31945次/s
3. Benchmark_Append-4 平均每秒执行的操作次数约为：1 / (29300 ns/operation) ≈ 34129次/s
4. Benchmark_Append-8 平均每秒执行的操作次数约为：1 / (24042 ns/operation) ≈ 41593次/s
5. Benchmark_Append-4 平均每秒执行的操作次数约为：1 / (30784 ns/operation) ≈ 32484次/s
6. Benchmark_Append-8 平均每秒执行的操作次数约为：1 / (30966 ns/operation) ≈ 32293次/s
7. Benchmark_Append-4 平均每秒执行的操作次数约为：1 / (35859 ns/operation) ≈ 27920次/s
8. Benchmark_Append-8 平均每秒执行的操作次数约为：1 / (33821 ns/operation) ≈ 29550次/s

- 获取数据性能估算

1. Benchmark_Get-4 平均每秒执行的操作次数约为：1 / (921 ns/operation) ≈  1085776次/s
2. Benchmark_Get-8 平均每秒执行的操作次数约为：1 / (636 ns/operation) ≈  1572327次/s
3. Benchmark_Get-4 平均每秒执行的操作次数约为：1 / (1558 ns/operation) ≈ 641848次/s
4. Benchmark_Get-8 平均每秒执行的操作次数约为：1 / (1296 ns/operation) ≈ 771604次/s
5. Benchmark_Get-4 平均每秒执行的操作次数约为：1 / (1695 ns/operation) ≈ 589970次/s
6. Benchmark_Get-8 平均每秒执行的操作次数约为：1 / (1402ns/operation) ≈  713266次/s
7. Benchmark_Get-4 平均每秒执行的操作次数约为：1 / (1865 ns/operation) ≈ 536000次/s
8. Benchmark_Get-8 平均每秒执行的操作次数约为：1 / (1730 ns/operation) ≈ 578034次/s

**写入数据性能**

- 在不同并发条件下，WFS存储引擎的写入操作平均每秒执行次数介于约 3万次/s 至 4万次/s 之间。

**读取数据性能**

- WFS存储引擎读数据操作的性能更为出色，平均每秒执行次数在 53万次/s 至 150万次/s 之间。

 **请注意：测试结果与环境有很大关系。实际应用中的性能可能会受到多种因素的影响，如系统负载、网络状况、磁盘I/O性能等，实际部署时需要根据具体环境进行验证和调优。**
 

------------

#### wfs内置图片基础处理

原图:   https://tlnet.top/statics/test/wfs_test.jpg
![](https://tlnet.top/statics/test/wfs_test.jpg)


- 裁剪正中部分，等比缩小生成200x200缩略图   https://tlnet.top/statics/test/wfs_test.jpg?imageView2/1/w/200/h/200
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/1/w/200/h/200)

- 宽度固定为200px，高度等比缩小，生成宽200缩略图    https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200)

- 高度固定为200px，宽度等比缩小，生成高200缩略图    https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/h/200
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/h/200)

- 高斯模糊，生成模糊程度Sigma为5，宽200的图片  https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/blur/5
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/blur/5)

- 灰色图片，生成灰色，宽200的图片   https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/grey/1
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/grey/1)

- 颜色反转，生成颜色相反，宽200的图片   https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/invert/1
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/invert/1)

- 水平反转 ，生成水平反转，宽200的图片   https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/fliph/1
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/fliph/1)

- 垂直反转 ，生成垂直反转，宽200的图片   https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/flipv/1
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/flipv/1)

- 图片旋转 ，生成向左旋转45度，宽200的图片   https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/rotate/45
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/rotate/45)

- 格式转换 ，生成向左旋转45，宽200的png图片   https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/rotate/45/format/png
![](https://tlnet.top/statics/test/wfs_test.jpg?imageView2/2/w/200/rotate/45/format/png)



##### 图片处理方式见 [wfs使用文档](https://tlnet.top/wfsdoc "wfs使用文档")

------------

#### WFS的使用简单说明

1. 执行文件下载地址：https://tlnet.top/download

2. 启动：
        ./linux101_wfs     -c    wfs.json

3.   wfs.json 配置说明

			{
   			 "listen": 4660,     
   			 "opaddr": ":6802",
    			"webaddr": ":6801",
   			 "memLimit": 128,
   	 		"data.maxsize": 10000,
  	 		 "filesize": 100,
			}
	
**属性说明：**

- listen                  http/https 资源获取服务监听端口
- opaddr               thrift后端资源操作地址
- webaddr            管理后台服务地址
- memLimit          wfs内存最大分配 (单位：MB)
- data.maxsize      wfs上传图片大小上限 (单位：KB)
- filesize                wfs后端归档文件大小上限 (单位：MB)

###### wfs使用详细说明请参考 [wfs使用文档](https://tlnet.top/wfsdoc "wfs使用文档")

------------

#### WFS如何存储，删除数据

1. http/https

		 curl -F "file=@1.jpg"  "http://127.0.0.1:6801/append/test/1.jpg" -H "username:admin" -H "password:123"

		 curl -X DELETE "http://127.0.0.1:6801/delete/test/1.jpg" -H "username:admin" -H "password:123"

2. 使用客户端

    以下是java客户端 示例

    	public void append() throws WfsException, IOException {
        String dir = System.getProperty("user.dir") + "/src/test/java/io/github/donnie4w/wfs/test/";
        WfsClient wc = newClient();
        WfsFile wf = new WfsFile();
        wf.setName("test/java/1.jpeg");
        wf.setData(Files.readAllBytes(Paths.get(dir + "1.jpeg")));
        wc.append(wf);
    	}

3. 通过管理后台上传/删除文件

------------

#### WFS管理后台

**默认搜索**
![](https://tlnet.top/f/1709440477_578.jpg)

**前缀搜索**
![](https://tlnet.top/f/1709440507_7665.jpg)

**碎片整理**
![](https://tlnet.top/f/1709440627_3436.jpg)

------------

#### WFS的分布式部署方案

wfs0.x版本到wfs1.x版本的设计变更说明：wfs0.x 版本实现了分布式存储，这使得系统能够跨多个服务器分散存储和处理数据，具备水平扩展能力和数据备份冗余能力，但是在实际应用中也暴露出一些问题，如元数据重复存储导致空间利用率不高。对于小文件的处理效率低，因为在节点间频繁转发传输，造成系统资源消耗增加。

wfs1.x版本的目标在于通过精简架构、聚焦性能提升来满足特定应用场景的需求，而在分布式部署方面的考量则交由用户借助第三方工具和服务来实现。
1. wfs1.x不直接支持分布式存储，但为了应对大规模部署和高可用需求，推荐采用如Nginx这样的负载均衡服务，通过合理的资源配置和定位策略，可以在逻辑上模拟出类似分布式的效果。也就是说，虽然每个wfs实例都是单机存储，但可以通过外部服务实现多个wfs实例之间的请求分发，从而达到业务层面的“分布式部署”。如何实现wfs的“分布式部署”可以参考文章《[WFS的分布式部署方案](https://tlnet.top/article/22425158 "WFS的分布式部署方案")》
2. 必须说明的是，超大规模数据存储业务中，分布式系统确实具有显著优势，包括动态资源调配、数据分块存储、多节点备份等高级功能。然而，分布式采用负载均衡策略的wfs1.x，则需要用户自行采取措施保证数据安全性和高可用性，例如定期备份数据、搭建负载均衡集群，并且在应用程序中配置并设计路由规则，确保数据能正确地路由到目标节点。
3. wfs的优势在于其简洁性和高效性。实际上，并非任何文件存储业务都需要复杂的分布式文件系统，相反，大部分业务尚未达到超大规模的量级，而使用复杂的分布式文件系统可能会带来与之不相匹配的过多额外成本和运维难度。目前的wfs及其相应的分布式部署策略已经能够较好地满足各种业务需求。
