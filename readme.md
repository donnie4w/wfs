# wfs是文件存储系统
主要是解决海量小文件存储的问题,服务器对海量小文件独立存储会出现许多问题

***

# 介绍
单个wfs可以单独运行 ，也可以多个wfs集群

#支持
 **上传文件，删除文件，拉取文件**
** 对图片文件输出大小处理：**
如：1.jpg?imageView2/0/w/100/h/100 输出宽高 100px的图片


# 启动wfs 
直接启动：**./wfs** （默认监听端口:3434） <br>
或带参数启动 ./wfs -max 50000000 -p 3434
**参数说明： -max是上传文件大小限制（单位字节）   -p启动端口（默认3434）** 
	


## 示例：

**curl 上传文件** <br/>

(1) 将文件直接上传服务器，地址为  http://*****:3434/u  ：如

	curl -F "file=@1.jpg" "http://127.0.0.1:3434/u"
 	则上传文件 1.jpg
	 上传完成后访问文件 ：http://127.0.0.1:3434/r/1.jpg 	
(2)将文件直接上传服务器，并指定文件名 地址为  http://*****:3434/u/文件名，
文件名同时也是访问的路径，如：

	curl -F "file=@1.jpg" "http://127.0.0.1:3434/u/abc/11"
	上传文件1.jpg 文件名 abc/11
	上传完成后访问文件 ：http://127.0.0.1:3434/r/abc/11

**curl 删除文件** 

	 curl -X DELETE "http://127.0.0.1:3434/d/1.jpg"
 	删除文件 1.jpg
 
	 curl -X DELETE "http://127.0.0.1:3434/d/abc/11"
	 删除文件 abc/11

***

## 支持在程序中使用客户端操作wfs

 	 wfsPost()    上传文件
 	 wfsRead()   拉取文件
 	 wfsDel()      删除文件

[以python客户端为例](https://github.com/donnie4w/wfs-pyclient "以python客户端为例")：

    	url ：http://*****:3434/thrift  固定形象
    	wfs = WfsClient("http://127.0.0.1:3434/thrift")
    	bs= getFileBytes("1.jpg")  获取图片
    	wfs.PostFile(bs,"aa/head.jpg","")   上传图片，并自定义图片路径aa/head.jpg
    	f=wfs.GetFile("aa/head.jpg") 拉取图片 aa/head.jpg 资源
    	print(len(f.fileBody))   
    	saveFileByBytes(f.fileBody,"22_1.jpg")
    	wfs.DelFile("aa/head.jpg")   删除图片
    	wfs.Close()

***

### **wfs的图片处理**
访问图片时，可以加参数来获取压缩后的图片 	<br/>
规则是：图片路径后+?imageView2/mode/w/Width/h/Height 如:

	imageView2固定
	mode 有 0，1，2，3，4，5 分别输出缩略图 限定的宽或高
	w后为宽：如：w/100
	h后为高：如：h/100
规则的制定是参考[七牛云存储](https://www.qiniu.com/ "七牛云存储")
所以mode规则也可以[参考](https://developer.qiniu.com/dora/api/1279/basic-processing-images-imageview2 "参考")

	http://127.0.0.1:3434/r/1.jpg?imageView2/0/w/100/h/100
	http://127.0.0.1:3434/r/1.jpg?imageView2/1/w/100/h/100 
	http://127.0.0.1:3434/r/1.jpg?imageView2/2/w/100
	http://127.0.0.1:3434/r/1.jpg?imageView2/3/h/100
	
	模式	 mode 说明
	/0/w/<LongEdge>/h/<ShortEdge>	限定缩略图的长边最多为<LongEdge>，短边最多为<ShortEdge>，进行等比缩放，不裁剪。如果只指定 w 参数则表示限定长边（短边自适应），只指定 h 参数则表示限定短边（长边自适应）。
	/1/w/<Width>/h/<Height>	限定缩略图的宽最少为<Width>，高最少为<Height>，进行等比缩放，居中裁剪。转后的缩略图通常恰好是 <Width>x<Height> 的大小（有一个边缩放的时候会因为超出矩形框而被裁剪掉多余部分）。如果只指定 w 参数或只指定 h 参数，代表限定为长宽相等的正方图。
	/2/w/<Width>/h/<Height>	限定缩略图的宽最多为<Width>，高最多为<Height>，进行等比缩放，不裁剪。如果只指定 w 参数则表示限定宽（高自适应），只指定 h 参数则表示限定高（宽自适应）。它和模式0类似，区别只是限定宽和高，不是限定长边和短边。从应用场景来说，模式0适合移动设备上做缩略图，模式2适合PC上做缩略图。
	/3/w/<Width>/h/<Height>	限定缩略图的宽最少为<Width>，高最少为<Height>，进行等比缩放，不裁剪。如果只指定 w 参数或只指定 h 参数，代表长宽限定为同样的值。你可以理解为模式1是模式3的结果再做居中裁剪得到的。
	/4/w/<LongEdge>/h/<ShortEdge>	限定缩略图的长边最少为<LongEdge>，短边最少为<ShortEdge>，进行等比缩放，不裁剪。如果只指定 w 参数或只指定 h 参数，表示长边短边限定为同样的值。这个模式很适合在手持设备做图片的全屏查看（把这里的长边短边分别设为手机屏幕的分辨率即可），生成的图片尺寸刚好充满整个屏幕（某一个边可能会超出屏幕）。
	/5/w/<LongEdge>/h/<ShortEdge>	限定缩略图的长边最少为<LongEdge>，短边最少为<ShortEdge>，进行等比缩放，居中裁剪。如果只指定 w 参数或只指定 h 参数，表示长边短边限定为同样的值。同上模式4，但超出限定的矩形部分会被裁剪。


***
**wfs提供了分片支持，分片可以解决单个节点的资源限制问题**

	wfs -slavelist 查询目前的节点 
	wfs -addslave slave1:192.168.1.101:3434  增加分片 节点名slave1，地址：192.168.1.101：3434
	wfs -addslave slave2:192.168.1.102:3434  增加分片 节点名slave2，地址：192.168.1.102：3434
	wfs -removeslave slave1  删除分片slave1

***

###  目前客户端有： python  java golang：
使用客户端操作 通过通讯协议的压缩传输 会更加快捷
1. [java : https://github.com/donnie4w/wfs-jclient](https://github.com/donnie4w/wfs-jclient)
2. [go : https://github.com/donnie4w/wfs-goclient](https://github.com/donnie4w/wfs-goclient)
3. [python : https://github.com/donnie4w/wfs-pyclient](https://github.com/donnie4w/wfs-pyclient)