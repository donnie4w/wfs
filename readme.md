# wfs是文件存储系统，主要是解决海量文件存储的问题,特别是小文件存储,原则上是简单易用,可扩展及备份恢复

***

# 介绍
单个wfs可以单独运行，多个wfs集群 可以启动wfs-slb  (github.com/donnie4w/wfs-slb) 作为代理层入口。
wfs没有过多额外功能，主要是 **增加文件，删除文件，拉取文件**


# 启动wfs 
**./wfs -max 50000000 -p 3434**
**参数说明： -max是上传文件大小限制（单位字节）   -p启动端口（默认3434）** 
	
使用wfs参考例子即可明白
# 1. 命令行
**上传文件** <br/>
(1)curl -F "file=@1.jpg" "http://127.0.0.1:3434/u"  <br/>
    **上传文件1.jpg 文件名 1.jpg** <br/>
(2)curl -F "file=@1.jpg" "http://127.0.0.1:3434/u/abc/11"   <br/>
    **上传文件1.jpg 文件名 abc/11** <br/>
例子(1)上传完成后访问文件 ：http://127.0.0.1:3434/r/1.jpg 	<br/>
例子(2)上传完成后访问文件 ：http://127.0.0.1:3434/r/abc/11   <br/>

**删除文件** <br/>
 curl -X DELETE "http://127.0.0.1:3434/d/1.jpg" <br/>    
 **删除文件 1.jpg**								<br/>
 curl -X DELETE "http://127.0.0.1:3434/d/abc/11"   <br/> 
 **删除文件 abc/11** 								<br/>

***

# 2. 使用thrift访问wfs     
  wfsPost()    上传文件
  wfsRead()    拉取文件  
  wfsDel       删除文件  
可以参考go版本  github.com/donnie4w/wfs-goclient  

***

wfs提供了一点附加的图片处理功能   	<br/>
访问图片时，可以加参数来获取压缩后的图片 	<br/>
参数规则与七牛图片的规则大致相同，（在本人多个项目中使用了七牛云存储，所以规则上希望能兼容七牛规则）	<br/>
https://developer.qiniu.com/dora/api/1279/basic-processing-images-imageview2	<br/>
imageView2/mode/w/Width/h/Height 
如： <br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/0/w/100/h/100 <br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/1/w/100/h/100 <br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/2/w/100	   <br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/3/h/100	   <br/>

mode 规则参考 https://developer.qiniu.com/dora/api/1279/basic-processing-images-imageview2 规则

***

分别打包了linux与windows两个执行文件	 <br/>
wfs-linux-amd64.gz		<br/>
wfs-windows-amd64.zip    <br/>
解压后 wfs --help 可以查看参数 ， 直接运行也可以默认端口3434  <br/>
