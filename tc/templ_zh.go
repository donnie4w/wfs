// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package tc

const (
	fileText = `<html>

    <head>
        <title>wfs</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link href="/bootstrap.css" rel="stylesheet">
        <script src="/bootstrap.min.js" type="text/javascript"></script>
    </head>
    
    <body class="container">
        <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand" href="/about">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-server"
                        viewBox="0 0 16 16">
                        <path
                            d="M1.333 2.667C1.333 1.194 4.318 0 8 0s6.667 1.194 6.667 2.667V4c0 1.473-2.985 2.667-6.667 2.667S1.333 5.473 1.333 4V2.667z" />
                        <path
                            d="M1.333 6.334v3C1.333 10.805 4.318 12 8 12s6.667-1.194 6.667-2.667V6.334a6.51 6.51 0 0 1-1.458.79C11.81 7.684 9.967 8 8 8c-1.966 0-3.809-.317-5.208-.876a6.508 6.508 0 0 1-1.458-.79z" />
                        <path
                            d="M14.667 11.668a6.51 6.51 0 0 1-1.458.789c-1.4.56-3.242.876-5.21.876-1.966 0-3.809-.316-5.208-.876a6.51 6.51 0 0 1-1.458-.79v1.666C1.333 14.806 4.318 16 8 16s6.667-1.194 6.667-2.667v-1.665z" />
                    </svg>
                </a>
                <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarCollapse">
                    <div class="navbar-nav">
                        <a class="nav-link" href='/init'>账号管理</a>
                        <a class="nav-link active" href='/file'>文件管理</a>
                        <a class="nav-link" href='/fragment'>碎片整理</a>
                        <a class="nav-link" href='/monitor'>系统监控</a>
                    </div>
                    <div class="navbar-nav ms-auto">
                        <a class="nav-link" href='/login'>登录</a>
                        <a class="nav-link" href="/lang?lang=en">[EN]</a>
                    </div>
                </div>
            </div>
        </nav>
        <div class="mt-1" style="font-size: xx-small;">
            <div class="container m-1" style="font-size: small;">
                <input id="lastId" type="text" hidden>
                <input id="searchType" type="text" hidden>
                <div class="input-group m-1">
                    <span class="input-group-text">默认搜索</span>
                    <select id="pagecount">
                        <option value="10">10</option>
                        <option value="50" selected>50</option>
                        <option value="100">100</option>
                        <option value="200">200</option>
                        <option value="500">500</option>
                    </select>
                    <input class="text-center" id="totalNum" placeholder="总数" style="width: 60px;" readonly />
                    <button class="btn btn-primary" onclick="search(-1);">上一页</button>
                    <input class="text-center" id="pageNumber" placeholder="页码" value="0" style="width: 50px;" />
                    <button class="btn btn-primary" onclick="search(1);">下一页</button>&nbsp;
                    <button class="btn btn-primary" onclick="search(0);">搜索</button>
                </div>
                <div class="input-group m-1">
                    <span class="input-group-text">前缀搜索</span>
                    <input id="prevName" placeholder="输入文件名前缀" style="width: 280px;" value="" />
                    <button class="btn btn-primary" onclick="search(2);">搜索</button>
                </div>
                <div class="input-group m-1">
                    <input class="btn btn-secondary" id="selectfile" type="button" value="选择文件" />
                    <span type="text" class="input-group-text" value="" id="showfilename"></span>
                    <input type="file" id="filebody" name="filebody" hidden />
                    <input id="filepath" placeholder="输入自定义访问路径" style="width: 230px;" value="" />
                    <input class="btn btn-primary" onclick="upload();" value="上传文件" type="button" />
                </div>
            </div>
            <table class="table table-striped" style="font-size: smaller;">
                <tr>
                    <th></th>
                    <th>文件名</th>
                    <th>大小(字节)</th>
                    <th>创建时间</th>
                    <th>图片预览</th>
                    <th>文件操作</th>
                </tr>
                <tbody id="fileTableBody">
                </tbody>
            </table>
            <div class="modal fade" id="renameModal" tabindex="-1" aria-labelledby="renameModalLabel" aria-hidden="true">
                <div class="modal-dialog">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h6 class="modal-title" id="renameModalLabel"></h6>
                        </div>
                        <div class="modal-body">
                            <input type="text" id="pathid" hidden>
                            <input type="text" class="form-control" id="newpath" placeholder="请输入新路径">
                        </div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">取消</button>
                            <button type="submit" class="btn btn-primary" onclick="renamesubmit()">确定</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <script>
            var aurl = "";
            var map = new Map();
            function search(t) {
                const formData = new FormData();
                const pageNumber = parseInt(document.getElementById("pageNumber").value);
                formData.append('pagecount', document.getElementById("pagecount").value);
                if (t == 0) {
                    formData.append('searchType', "1");
                    formData.append('pageNumber', pageNumber);
                } else if (t == -1) {
                    formData.append('searchType', "1");
                    formData.append('pageNumber', pageNumber - 1);
                } else if (t == 1) {
                    formData.append('searchType', "1");
                    formData.append('pageNumber', pageNumber + 1);
                    formData.append('lastId', document.getElementById("lastId").value);
                    document.getElementById("lastId").value = ""
                } else if (t == 2) {
                    const prevName = document.getElementById("prevName").value
                    if (prevName == "") {
                        return
                    }
                    formData.append('prevName', document.getElementById("prevName").value);
                    formData.append('searchType', "2");
                }
    
                fetch('/filedata', {
                    method: 'POST',
                    body: formData,
                }).then(response => {
                    if (!response.ok) {
                        throw new Error("HTTP error! status: ${response.status}");
                    }
                    return response.json();
                }).then(data => {
                    document.getElementById("fileTableBody").innerHTML = ""
                    const RevProxy = data.RevProxy;
                    const ClientPort = data.ClientPort;
                    const CliProtocol = data.CliProtocol
                    document.getElementById("totalNum").value = data.TotalNum;
                    document.getElementById("pageNumber").value = data.CurrentNum;
                    const url = "/r/"
    
                    if (aurl == "") {
                        aurl = CliProtocol + window.location.hostname + ":" + ClientPort + "/"
                        if (RevProxy != "") {
                            aurl = RevProxy
                        }
                    }
                    map.clear()
                    let id = 1;
                    const fs = data.FS
                    if (Array.isArray(fs) && fs.length > 0) {
                        for (const item of fs) {
                            map.set(item.Id, item.Name);
                            let tr = document.createElement('tr');
                            let d = '<td style="font-weight: bold;">' + item.Id + '</td>'
                                + '<td id="' + item.Id + '">' + item.Name + '</td>'
                                + '<td>' + item.Size + '</td>'
                                + '<td>' + item.Time + '</td>'
                                + '<td><a id="a' + item.Id + '" href="' + aurl + item.Name + '" target="_blank"><img src="' + url + item.Name + "?mode/2/h/60/" + Date.now() + '" height="60" onerror="this.src=\'data:image/svg+xml;charset=utf-8,\' + encodeURIComponent(svgCode)" alt="Fallback Image" /></a></td>'
                                + '<td><button class="btn btn-primary m-1" onclick=\'renameshow(' + item.Id + ',"' + item.Name + '")\';">重命名</button><button class="btn btn-primary" onclick=\'deletefile(this,"' + item.Name + '")\';">删除</button></td>'
                            tr.innerHTML = d;
                            document.getElementById("fileTableBody").appendChild(tr);
                            document.getElementById("lastId").value = item.Id
                        }
                    }
                }).catch(error => {
                    console.error('Error submitting form:', error);
                });
            }
    
            document.getElementById('selectfile').addEventListener('click', function () {
                clearFileInput();
                document.getElementById('filebody').click();
            });
    
            document.getElementById('filebody').addEventListener('change', function () {
                document.getElementById("showfilename").innerHTML = document.getElementById('filebody').files[0].name + "<a href='javascript:;' onclick='clearFileInput()'>&#10006;</a>";
            });
    
            function clearFileInput() {
                document.getElementById("showfilename").innerHTML = ""
                document.getElementById('filebody').value = ""
            }
    
            function upload() {
                const filebody = document.getElementById("filebody");
                if (filebody.files.length == 0) {
                    return
                }
                const file = filebody.files[0];
                getLimit().then(data => {
                    let limit = data.limit;
                    if (limit > 0 && file.size > limit) {
                        if (limit > 1 << 20) {
                            alert("data oversize \n\nthe maximum data size is " + limit / (1 << 20) + "MB")
                        } else {
                            alert("data oversize \n\nthe maximum data size is " + limit / (1 << 10) + "KB")
                        }
                        return
                    }
    
                    const formData = new FormData();
                    formData.append("file", file, file.name)
                    const filepath = document.getElementById("filepath").value
                    if (filepath != "") {
                        formData.append("filename", filepath)
                    }
    
                    fetch('/append/', {
                        method: 'POST',
                        body: formData,
                    }).then(response => {
                        if (!response.ok) {
                            throw new Error("HTTP error! status: ${response.status}");
                        }
                        return response.json();
                    }).then(data => {
                        if (data.status) {
                            alert("上传文件成功:" + data.name)
                            clearFileInput();
                            document.getElementById("filepath").value = "";
                        } else {
                            alert("上传文件失败：" + data.desc)
                        }
                    })
    
                })
    
            }
    
            async function getLimit() {
                const formData = new FormData();
                formData.append("limit", "0")
                let response = await fetch('/file', {
                    method: 'POST',
                    body: formData,
                })
                return await response.json();
            }
    
            function deletefile(t, name) {
                const formData = new FormData();
                formData.append("filename", name)
                if (confirm("确定删除该文件？")) {
                    fetch('/delete/', {
                        method: 'DELETE',
                        body: formData,
                    }).then(response => {
                        if (!response.ok) {
                            throw new Error("HTTP error! status: ${response.status}");
                        }
                        return response.json();
                    }).then(data => {
                        if (data.status) {
                            console.log('Delete successful:', data);
                            document.getElementById("fileTableBody").removeChild(t.parentNode.parentNode);
                        } else {
                            alert("删除文件失败：" + data.desc)
                        }
                    }).catch(error => {
                        console.error('Error:', error);
                    });
                }
            }
    
            const rnModal = new bootstrap.Modal(document.getElementById('renameModal'));
    
            function renameshow(id) {
                rnModal.show();
                document.getElementById("pathid").value = id;
                document.getElementById("renameModalLabel").innerText = document.getElementById(id).innerText;
                document.getElementById("newpath").value = document.getElementById(id).innerText;
            }
    
            function renameModalhide() {
                document.getElementById("pathid").value = "";
                document.getElementById("newpath").value = "";
                document.getElementById("renameModalLabel").innerText = "";
                rnModal.hide();
            }
    
            function renamesubmit() {
                const formData = new FormData();
                formData.append("path", map.get(parseInt(document.getElementById("pathid").value)))
                formData.append("newpath", document.getElementById("newpath").value)
                fetch('/rename', {
                    method: 'POST',
                    body: formData,
                }).then(response => {
                    if (!response.ok) {
                        throw new Error("HTTP error! status: ${response.status}");
                    }
                    return response.json();
                }).then(data => {
                    if (data.status) {
                        let id = document.getElementById('pathid').value;
                        document.getElementById(id).innerText = document.getElementById("newpath").value;
                        document.getElementById("a" + id).href = aurl + document.getElementById("newpath").value;
                        map.set(parseInt(id), document.getElementById("newpath").value)
                        renameModalhide();
                    } else {
                        alert("重命名失败:" + data.desc)
                    }
                }).catch(error => {
                    console.error('Error:', error);
                });
            }
        </script>
        <script>
            const svgCode = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-file-earmark-text" viewBox="0 0 16 16">'
                + '<path d="M5.5 7a.5.5 0 0 0 0 1h5a.5.5 0 0 0 0-1h-5zM5 9.5a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5zm0 2a.5.5 0 0 1 .5-.5h2a.5.5 0 0 1 0 1h-2a.5.5 0 0 1-.5-.5z"/>'
                + '<path d="M9.5 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V4.5L9.5 0zm0 1v2A1.5 1.5 0 0 0 11 4.5h2V14a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1h5.5z"/>'
                + '</svg>';
        </script>
    </body>
    
    </html>`

	fragmentText = `
    <html>

<head>
    <title>wfs</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
    <script src="/bootstrap.min.js" type="text/javascript"></script>
</head>

<body class="container">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
            <a class="navbar-brand" href="/about">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-server"
                    viewBox="0 0 16 16">
                    <path
                        d="M1.333 2.667C1.333 1.194 4.318 0 8 0s6.667 1.194 6.667 2.667V4c0 1.473-2.985 2.667-6.667 2.667S1.333 5.473 1.333 4V2.667z" />
                    <path
                        d="M1.333 6.334v3C1.333 10.805 4.318 12 8 12s6.667-1.194 6.667-2.667V6.334a6.51 6.51 0 0 1-1.458.79C11.81 7.684 9.967 8 8 8c-1.966 0-3.809-.317-5.208-.876a6.508 6.508 0 0 1-1.458-.79z" />
                    <path
                        d="M14.667 11.668a6.51 6.51 0 0 1-1.458.789c-1.4.56-3.242.876-5.21.876-1.966 0-3.809-.316-5.208-.876a6.51 6.51 0 0 1-1.458-.79v1.666C1.333 14.806 4.318 16 8 16s6.667-1.194 6.667-2.667v-1.665z" />
                </svg>
            </a>
            <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">
                <div class="navbar-nav">
                    <a class="nav-link" href='/init'>账号管理</a>
                    <a class="nav-link" href='/file'>文件管理</a>
                    <a class="nav-link active" href='/fragment'>碎片整理</a>
                    <a class="nav-link" href='/monitor'>系统监控</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>登录</a>
                    <a class="nav-link" href="/lang?lang=en">[EN]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: xx-small;">
        <table class="table table-striped" style="font-size: smaller;">
            <tr>
                <th>文件名</th>
                <th>最后更新时间</th>
                <th>文件大小(字节)</th>
                <th>碎片大小(字节)</th>
                <th>文件状态</th>
                <th>文件操作</th>
            </tr>
            <tbody id="fileTableBody">
                {{range $k,$v := .}}
                <tr>
                    <td>{{$v.Name}}</td>
                    <td>{{$v.Time}}</td>
                    <td>{{$v.FileSize}}</td>
                    <td>{{$v.FragmentSize}}</td>
                    {{if eq $v.Status 1 }}
                    <td>只读</td>
                    {{else}}
                    <td>读写中(禁止操作)</td>
                    {{end}}
                    {{if and (eq $v.Status 1) (gt $v.FragmentSize 0)}}
                    <td><button id="{{$v.Name}}" class="btn btn-primary" onclick="fragment('{{$v.Name}}')">碎片整理</button>
                    </td>
                    {{else}}
                    <td></td>
                    {{end}}
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
    <script>
        function fragment(n) {
            if (confirm("[" + n + "]确定进行碎片整理？\n建议碎片整理应当在WFS服务操作比较少的状态进行，可减少对前端服务质量的影响")) {
                document.getElementById(n).innerText = "碎片整理中..."
                document.getElementById(n).disabled = true;
                const formData = new FormData();
                formData.append("node", n)
                fetch('/defrag', {
                    method: 'POST',
                    body: formData,
                }).then(response => {
                    if (!response.ok) {
                        throw new Error("HTTP error! status: ${response.status}");
                    }
                    document.getElementById(n).innerText = "碎片整理"
                    document.getElementById(n).disabled = false;
                    return response.json();
                }).then(data => {
                    if (data.status) {
                        window.location.reload(true);
                    } else {
                        alert("碎片整理失败")
                    }
                }).catch(error => {
                    console.error('Error:', error);
                });
            }
        }
    </script>
</body>

</html>
    `

	initText = `
    <html>

<head>
    <title>wfs</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
    <script src="/bootstrap.min.js" type="text/javascript"></script>
</head>

<body class="container">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
            <a class="navbar-brand" href="/about">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-server"
                    viewBox="0 0 16 16">
                    <path
                        d="M1.333 2.667C1.333 1.194 4.318 0 8 0s6.667 1.194 6.667 2.667V4c0 1.473-2.985 2.667-6.667 2.667S1.333 5.473 1.333 4V2.667z" />
                    <path
                        d="M1.333 6.334v3C1.333 10.805 4.318 12 8 12s6.667-1.194 6.667-2.667V6.334a6.51 6.51 0 0 1-1.458.79C11.81 7.684 9.967 8 8 8c-1.966 0-3.809-.317-5.208-.876a6.508 6.508 0 0 1-1.458-.79z" />
                    <path
                        d="M14.667 11.668a6.51 6.51 0 0 1-1.458.789c-1.4.56-3.242.876-5.21.876-1.966 0-3.809-.316-5.208-.876a6.51 6.51 0 0 1-1.458-.79v1.666C1.333 14.806 4.318 16 8 16s6.667-1.194 6.667-2.667v-1.665z" />
                </svg>
            </a>
            <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">
                <div class="navbar-nav">
                    <a class="nav-link active" href='/init'>账号管理</a>
                    <a class="nav-link" href='/file'>文件管理</a>
                    <a class="nav-link" href='/fragment'>碎片整理</a>
                    <a class="nav-link" href='/monitor'>系统监控</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>登录</a>
                    <a class="nav-link" href="/lang?lang=en">[EN]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: small;">
        {{if .ShowCreate }}
        <div class="container-fluid card mt-1 p-1">
            </h6>
            <form class="form-control" id="createAdminform" action="/init?type=1" method="post">
                <h6>新建管理员 <h6 class="important">{{ .Show }}</h6>
                    <input name="adminName" placeholder="用户名" />
                    <input name="adminPwd" placeholder="密码" type="password" />
                    管理员<input name="adminType" type="radio" value="1" checked />
                    {{if not .Init}}
                    观察员<input name="adminType" type="radio" value="2" />
                    {{end}}
                    <input type="submit" class="btn btn-primary" value="新建管理员" />
            </form>
        </div>
        {{end}}
        {{if not .Init}}
        <div class="container-fluid card mt-1 p-1">
            <div class="m-2">
                <h6>后台管理员</h6>
                {{range $k,$v := .AdminUser}}
                <form class="form-control" id="adminform" action="/init?type=2" method="post">
                    <input name="adminName" value='{{ $k }}' readonly style="border:none;" /> 权限:{{ $v }}
                    <input type="button" value="删除用户" class="btn btn-danger"
                        onclick="javascipt:if (confirm('确定删除?')){this.parentNode.submit();};" />
                </form>
                {{end}}
            </div>
        </div>
        <hr>
        {{end}}
    </div>

</html>
    `

	loginText = `
    <html>

<head>
    <title>wfs</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
</head>

<body class="container">
    <div class="container-fluid text-right">
        <span>
            <h4 style="display: inline;">WFS 管理后台</h4>
        </span>
        <span style="text-align:right">
            <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
        </span>
        <hr>
        <div id="login">
            <h5>登录</h5>
            <form class="form-control" id="loginform" action="/login" method="post">
                <input name="type" value="1" hidden />
                <input name="name" placeholder="用户名" />
                <input name="pwd" placeholder="密码" type="password" />
                <input type="submit" class="btn btn-primary" value="登录" />
            </form>
        </div>
        <hr>
    </div>
</html>
    `

	monitorText = `
    <html>

<head>
    <title>wfs</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
    <script src="/bootstrap.min.js" type="text/javascript"></script>
</head>

<body class="container">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
            <a class="navbar-brand" href="/about">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-server"
                    viewBox="0 0 16 16">
                    <path
                        d="M1.333 2.667C1.333 1.194 4.318 0 8 0s6.667 1.194 6.667 2.667V4c0 1.473-2.985 2.667-6.667 2.667S1.333 5.473 1.333 4V2.667z" />
                    <path
                        d="M1.333 6.334v3C1.333 10.805 4.318 12 8 12s6.667-1.194 6.667-2.667V6.334a6.51 6.51 0 0 1-1.458.79C11.81 7.684 9.967 8 8 8c-1.966 0-3.809-.317-5.208-.876a6.508 6.508 0 0 1-1.458-.79z" />
                    <path
                        d="M14.667 11.668a6.51 6.51 0 0 1-1.458.789c-1.4.56-3.242.876-5.21.876-1.966 0-3.809-.316-5.208-.876a6.51 6.51 0 0 1-1.458-.79v1.666C1.333 14.806 4.318 16 8 16s6.667-1.194 6.667-2.667v-1.665z" />
                </svg>
            </a>
            <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">

                <div class="navbar-nav">
                    <a class="nav-link" href='/init'>账号管理</a>
                    <a class="nav-link" href='/file'>文件管理</a>
                    <a class="nav-link" href='/fragment'>碎片整理</a>
                    <a class="nav-link active" href='/monitor'>系统监控</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>登录</a>
                    <a class="nav-link" href="/lang?lang=en">[EN]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="container mt-1 card" style="font-size: small;">
        <div class="container mt-1" style="font-size: small;">
            <h3>性能数据监控</h3>
            <div class="input-group">
                <span class="input-group-text">监控时间间隔(单位:秒)</span>
                <input id="stime" placeholder="输入时间间隔" value="3" />
                <button class="btn btn-primary" onclick="monitorLoad();">开始</button>&nbsp;
                <button class="btn btn-primary" onclick="stop();">停止</button>&nbsp;
                <button class="btn btn-primary" onclick="clearData();">清除数据</button>
            </div>
        </div>

        <table class="table table-striped " style="font-size: smaller;">
            <tr>
                <th></th>
                <th>内存使用(MB)</th>
                <th>内存已分配(MB)</th>
                <th>内存释放次数</th>
                <th>协程数</th>
                <th>CPU核数</th>
                <th>磁盘剩余(GB)</th>
                <th>内存使用率</th>
                <th>CPU使用率</th>
            </tr>
            <tbody id="monitorBody">
            </tbody>
        </table>
    </div>
</body>
<script type="text/javascript">
    var pro = window.location.protocol;
    var wspro = "ws:";
    if (pro === "https:") {
        wspro = "wss:";
    }
    var wsmnt = null;
    var id = 1;
    function WS() {
        this.ws = null;
    }

    WS.prototype.monitor = function () {
        let obj = this;
        this.ws = new WebSocket(wspro + "//" + window.location.host + "/monitorData");
        this.ws.onopen = function (evt) {
            obj.ws.send(document.getElementById("stime").value);
        }
        this.ws.onmessage = function (evt) {
            if (evt.data != "") {
                var json = JSON.parse(evt.data);
                var tr = document.createElement('tr');
                var d = '<td style="font-weight: bold;">' + id++ + '</td>'
                    + '<td>' + Math.round(json.Alloc / (1 << 20)) + '</td>'
                    + '<td>' + Math.round(json.TotalAlloc / (1 << 20)) + '</td>'
                    + '<td>' + json.NumGC + '</td>'
                    + '<td>' + json.NumGoroutine + '</td>'
                    + '<td>' + json.NumCPU + '</td>'
                    + '<td>' + json.DiskFree + '</td>'
                    + '<td>' + Math.round(json.RamUsage * 10000) / 100 + '%</td>'
                    + '<td>' + Math.round(json.CpuUsage * 100) / 100 + '%</td>';
                tr.innerHTML = d;
                document.getElementById("monitorBody").appendChild(tr);
            }
        }
    }

    WS.prototype.close = function () {
        this.ws.close();
    }

    function monitorLoad() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
        wsmnt = new (WS);
        wsmnt.monitor();
    }

    function stop() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
    }

    function clearData() {
        document.getElementById("monitorBody").innerHTML = "";
    }

</script>

</html>
    `

    aboutText=`<html>

    <head>
        <title>wfs</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link href="/bootstrap.css" rel="stylesheet">
        <script src="/bootstrap.min.js" type="text/javascript"></script>
    </head>
    
    <body class="container">
        <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand" href="/about"">
                    <svg xmlns=" http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                    class="bi bi-server" viewBox="0 0 16 16">
                    <path
                        d="M1.333 2.667C1.333 1.194 4.318 0 8 0s6.667 1.194 6.667 2.667V4c0 1.473-2.985 2.667-6.667 2.667S1.333 5.473 1.333 4V2.667z" />
                    <path
                        d="M1.333 6.334v3C1.333 10.805 4.318 12 8 12s6.667-1.194 6.667-2.667V6.334a6.51 6.51 0 0 1-1.458.79C11.81 7.684 9.967 8 8 8c-1.966 0-3.809-.317-5.208-.876a6.508 6.508 0 0 1-1.458-.79z" />
                    <path
                        d="M14.667 11.668a6.51 6.51 0 0 1-1.458.789c-1.4.56-3.242.876-5.21.876-1.966 0-3.809-.316-5.208-.876a6.51 6.51 0 0 1-1.458-.79v1.666C1.333 14.806 4.318 16 8 16s6.667-1.194 6.667-2.667v-1.665z" />
                    </svg>
                </a>
                <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarCollapse">
                    <div class="navbar-nav">
                        <a class="nav-link" href='/init'>账号管理</a>
                        <a class="nav-link" href='/file'>文件管理</a>
                        <a class="nav-link" href='/fragment'>碎片整理</a>
                        <a class="nav-link" href='/monitor'>系统监控</a>
                    </div>
                    <div class="navbar-nav ms-auto">
                        <a class="nav-link" href='/login'>登录</a>
                        <a class="nav-link" href="/lang?lang=en">[EN]</a>
                    </div>
                </div>
            </div>
        </nav>
        <div class="container my-1 card">
            <h3>WFS 简介 </h3>
            <h6>版权所有, 作者 <a href="mailto:donnie4w@gmail.com" class="text-reset">donnie4w</a>, 版本 v{{.}} </h6>
            <hr>
            <div class="my-2">
                <p><strong>WFS 文件存储系统</strong>，主要用于解决海量小文件的存储问题。</p>
                <p>wfs官网：<a href="https://tlnet.top/wfs" target="_blank">https://tlnet.top/wfs</a></p>
                <p><a href="https://tlnet.top/wfsdoc" target="_blank">wfs使用文档</a></p>
            </div>
            <hr>
            <h5 class="mt-2">wfs 相关程序</h5>
            <div class="row button-links">
                <div class="col-sm-6 col-md-4 col-lg-3 mb-3">
                    <a href="https://github.com/donnie4w/wfs" target="_blank" class="btn btn-primary">wfs 源码地址</a>
                </div>
                <div class="col-sm-6 col-md-4 col-lg-3 mb-3">
                    <a href="https://github.com/donnie4w/wfs-goclient" target="_blank" class="btn btn-secondary">Go 客户端</a>
                </div>
                <div class="col-sm-6 col-md-4 col-lg-3 mb-3">
                    <a href="https://github.com/donnie4w/wfs-jclient" target="_blank" class="btn btn-secondary">Java 客户端</a>
                </div>
                <div class="col-sm-6 col-md-4 col-lg-3 mb-3">
                    <a href="https://github.com/donnie4w/wfs-pyclient" target="_blank" class="btn btn-secondary">Python
                        客户端</a>
                </div>
            </div>
    
            <p class="text-muted mt-3">
                <strong>Email:</strong>
                <a href="mailto:donnie4w@gmail.com" class="text-reset">donnie4w@gmail.com</a>
            </p>
        </div>
    
    </html>`

)
