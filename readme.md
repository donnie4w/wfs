wfs是文件存储系统，主要是解决海量文件存储的问题，特别是小文件存储 <br/>
原则上是简单易用，提供可扩展，备份恢复等功能 <br/>
<br/>
单个wfs可以单独运行，多个wfs集群 可以启动wfs-slb  (github.com/donnie4w/wfs-slb) 作为代理层入口  <br/>
wfs使用比较简单
启动wfs 
    ./wfs -max 50000000 -p 3434  <br/>
	 -max 是上传文件大小限制 单位字节 <br/>
	 -p 端口 （默认3434） <br/>
	
使用wfs参考例子即可明白 <br/>
1.<br/>
上传文件：<br/>
 (1)curl -F "file=@1.jpg" "http://127.0.0.1:3434/u"           // 上传文件1.jpg 文件名 1.jpg <br/>
 (2)curl -F "file=@1.jpg" "http://127.0.0.1:3434/u/abc/11"    // 上传文件1.jpg 文件名 abc/11 <br/>

  例子(1)上传完成后访问文件 ：http://127.0.0.1:3434/r/1.jpg <br/>
  例子(2)上传完成后访问文件 ：http://127.0.0.1:3434/r/abc/11 <br/>

删除文件 <br/>
 curl -X DELETE "http://127.0.0.1:3434/d/1.jpg"     //删除文件 1.jpg <br/>
 curl -X DELETE "http://127.0.0.1:3434/d/abc/11"    //删除文件 abc/11 <br/>

2.使用thrift访问wfs     <br/>
  wfsPost()    //上传文件  <br/>
  wfsRead()    //拉取文件  <br/>
  wfsDel       //删除文件  <br/>
可以参考go版本  github.com/donnie4w/wfs-goclient  <br/>

wfs提供了一点附加的图片处理功能   	<br/>
访问图片时，可以加参数来获取压缩后的图片 	<br/>
参数规则与七牛图片的规则大致相同，（在本人多个项目中使用了七牛云存储，所以规则上希望能兼容七牛规则）	<br/>
https://developer.qiniu.com/dora/api/1279/basic-processing-images-imageview2	<br/>
imageView2/mode/w/LongEdge <br/>
                 /h/ShortEdge <br/>
                 /format/Format	<br/>
                 /interlace/Interlace	<br/>
                 /q/Quality			<br/>
                 /ignore-error/ignoreError	<br/>
mode 规则 <br/>
/0/w/LongEdge>/h/<ShortEdge 限定缩略图的长边最多为LongEdge，短边最多为ShortEdge，进行等比缩放，不裁剪。如果只指定 w 参数则表示限定长边（短边自适应），只指定 h 参数则表示限定短边（长边自适应）。<br/>
/1/w/Width/h/Height	限定缩略图的宽最少为Width，高最少为Height，进行等比缩放，居中裁剪。转后的缩略图通常恰好是 Width x Height 的大小（有一个边缩放的时候会因为超出矩形框而被裁剪掉多余部分）。如果只指定 w 参数或只指定 h 参数，代表限定为长宽相等的正方图。<br/>
/2/w/Width/h/Height	限定缩略图的宽最多为Width，高最多为Height，进行等比缩放，不裁剪。如果只指定 w 参数则表示限定宽（长自适应），只指定 h 参数则表示限定长（宽自适应）。它和模式0类似，区别只是限定宽和高，不是限定长边和短边。从应用场景来说，模式0适合移动设备上做缩略图，模式2适合PC上做缩略图。<br/>
/3/w/Width/h/Height	限定缩略图的宽最少为Width，高最少为Height，进行等比缩放，不裁剪。如果只指定 w 参数或只指定 h 参数，代表长宽限定为同样的值。你可以理解为模式1是模式3的结果再做居中裁剪得到的。<br/>
/4/w/LongEdge/h/ShortEdge 限定缩略图的长边最少为LongEdge，短边最少为ShortEdge，进行等比缩放，不裁剪。如果只指定 w 参数或只指定 h 参数，表示长边短边限定为同样的值。这个模式很适合在手持设备做图片的全屏查看（把这里的长边短边分别设为手机屏幕的分辨率即可），生成的图片尺寸刚好充满整个屏幕（某一个边可能会超出屏幕）。 <br/>
/5/w/LongEdge/h/ShortEdge 限定缩略图的长边最少为LongEdge，短边最少为ShortEdge，进行等比缩放，居中裁剪。如果只指定 w 参数或只指定 h 参数，表示长边短边限定为同样的值。同上模式4，但超出限定的矩形部分会被裁剪。 <br/>				

如： <br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/0/w/100/h/100 <br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/1/w/100/h/100 <br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/2/w/100	<br/>
http://127.0.0.1:3434/r/1.jpg?imageView2/3/h/100	<br/>

分别打包了linux与windows两个执行文件	 <br/>
wfs-linux-amd64.gz		<br/>
wfs-windows-amd64.zip    <br/>
解压后 wfs --help 可以查看参数 ， 直接运行也可以默认端口3434  <br/>
