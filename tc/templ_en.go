// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package tc

const (
	loginEnText = `
    <html>

<head>
    <title>wfs</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
</head>

<body class="container">
    <div class="container-fluid text-right">
        <span>
            <h4 style="display: inline;">WFS Management Platform</h4>
        </span>
        <span style="text-align:right">
            <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=zh">[中文]</a></h6>
        </span>
        <hr>
        <div id="login">
            <h5>Login</h5>
            <form class="form-control" id="loginform" action="/login" method="post">
                <input name="type" value="1" hidden />
                <input name="name" placeholder="username" />
                <input name="pwd" placeholder="password" type="password" />
                <input type="submit" class="btn btn-primary" value="Login" />
            </form>
        </div>
        <hr>
    </div>
</html>
    `

	initEnText = `
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
                    <a class="nav-link active" href='/init'>Account Management</a>
                    <a class="nav-link" href='/file'>File Management</a>
                    <a class="nav-link" href='/fragment'>Fragmentation Cleanup</a>
                    <a class="nav-link" href='/monitor'>System Monitoring</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>Log In</a>
                    <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: small;">
        {{if .ShowCreate }}
        <div class="container-fluid card mt-1 p-1">
            </h6>
            <form class="form-control" id="createAdminform" action="/init?type=1" method="post">
                <h6>Create Administrator <h6 class="important">{{ .Show }}</h6>
                    <input name="adminName" placeholder="username" />
                    <input name="adminPwd" placeholder="password" type="password" />
                    Administrator<input name="adminType" type="radio" value="1" checked />
                    {{if not .Init}}
                    Observer<input name="adminType" type="radio" value="2" />
                    {{end}}
                    <input type="submit" class="btn btn-primary" value="Create" />
            </form>
        </div>
        {{end}}
        {{if not .Init}}
        <div class="container-fluid card mt-1 p-1">
            <div class="m-2">
                <h6>Manage Accounts</h6>
                {{range $k,$v := .AdminUser}}
                <form class="form-control" id="adminform" action="/init?type=2" method="post">
                    <input name="adminName" value='{{ $k }}' readonly style="border:none;" /> Authority:{{ $v }}
                    <input type="button" value="Delete Account" class="btn btn-danger"
                        onclick="javascipt:if (confirm('confirm delete?')){this.parentNode.submit();};" />
                </form>
                {{end}}
            </div>
        </div>
        <hr>
        {{end}}
    </div>

</html>
    `

	fileEnText = `<html>

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
                        d="M1.333 2.667C1.333 1.194 4.318 0 8 0s6.667 1.194 6.667 2.667V4c0 1.473-2.985 2.667-6.667 2.667S1.333 5.473 1.333 4V2.667z"/>
                <path
                        d="M1.333 6.334v3C1.333 10.805 4.318 12 8 12s6.667-1.194 6.667-2.667V6.334a6.51 6.51 0 0 1-1.458.79C11.81 7.684 9.967 8 8 8c-1.966 0-3.809-.317-5.208-.876a6.508 6.508 0 0 1-1.458-.79z"/>
                <path
                        d="M14.667 11.668a6.51 6.51 0 0 1-1.458.789c-1.4.56-3.242.876-5.21.876-1.966 0-3.809-.316-5.208-.876a6.51 6.51 0 0 1-1.458-.79v1.666C1.333 14.806 4.318 16 8 16s6.667-1.194 6.667-2.667v-1.665z"/>
            </svg>
        </a>
        <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarCollapse">
            <div class="navbar-nav">
                <a class="nav-link" href='/init'>Account Management</a>
                <a class="nav-link active" href='/file'>File Management</a>
                <a class="nav-link" href='/fragment'>Fragmentation Cleanup</a>
                <a class="nav-link" href='/monitor'>System Monitoring</a>
            </div>
            <div class="navbar-nav ms-auto">
                <a class="nav-link" href='/login'>Log In</a>
                <a class="nav-link" href="/lang?lang=zh">[中文]</a>
            </div>
        </div>
    </div>
</nav>
<div class="mt-1" style="font-size: xx-small;">
    <div class="container m-1" style="font-size: small;">
        <div class="input-group m-1">
            <input id="lastId" type="text" hidden>
            <input id="searchType" type="text" hidden>
            <span class="input-group-text">Default Search</span>
            <select id="pagecount">
                <option value="10">10</option>
                <option value="50" selected>50</option>
                <option value="100">100</option>
                <option value="200">200</option>
                <option value="500">500</option>
            </select>
            <input class="text-center" id="totalNum" placeholder="Total" style="width: 60px;" readonly/>
            <button class="btn btn-primary" onclick="search(-1);">Previous</button>
            <input class="text-center" id="pageNumber" placeholder="page number" value="0" style="width: 50px;"/>
            <button class="btn btn-primary" onclick="search(1);">Next</button>&nbsp;
            <button class="btn btn-primary" onclick="search(0);">Search</button>
        </div>
        <div class="input-group m-1">
            <span class="input-group-text">Search by Prefix</span>
            <input id="prevName" placeholder="Enter file prefix" style="width: 280px;" value=""/>
            <button class="btn btn-primary" onclick="search(2);">Search</button>
        </div>
        <div class="input-group m-1">
            <input class="btn btn-secondary" id="selectfile" type="button" value="Select File"/>
            <span type="text" class="input-group-text" value="" id="showfilename"></span>
            <input type="file" id="filebody" name="filebody" hidden/>
            <input id="filepath" placeholder="Enter Custom Access Path" style="width: 230px;" value=""/>
            <input class="btn btn-primary" onclick="upload();" value="Upload File" type="button"/>
        </div>
    </div>
    <table class="table table-striped" style="font-size: smaller;">
        <tr>
            <th></th>
            <th>File Name</th>
            <th>Size (Bytes)</th>
            <th>Create Time</th>
            <th>Image Preview</th>
            <th>File Operations</th>
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
                    <input type="text" class="form-control" id="newpath" placeholder="Enter a new path">
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">cancel</button>
                    <button type="submit" class="btn btn-primary" onclick="renamesubmit()">submit</button>
                </div>
            </div>
        </div>
    </div>

    <div class="modal fade" id="processpicModal" tabindex="-1" aria-labelledby="processpicModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-lg">
            <div class="modal-content">
                <div class="modal-header">
                    <h6 class="modal-title" id="processpicModalLabel"></h6>
                    <button type="button" class="btn btn-primary" onclick="copysrc()">复制</button>
                </div>
                <div class="modal-body">
                    <div class="row g-2">
                        <div class="col-8">
                            <img src="" id="picsrcid" class="img-fluid mx-auto d-block" style="max-width: 100%; height: auto; object-fit: contain; width: auto;">
                        </div>
                        <div class="col-2">
                            <div class="d-grid gap-2">
                                <button type="button" class="btn btn-secondary btn-sm" onclick="rotateLeft()">L Rot</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="zoomIn()">Zoom In</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="fliph()">H-Flip</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="invert()">Invert</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picblur(0.5)">Blur</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picformat('jpeg')">JPEG</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picformat('ico')">ICO</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picformat('bmp')">BMP</button>
                            </div>
                        </div>
                        <div class="col-2">
                            <div class="d-grid gap-2">
                                <button type="button" class="btn btn-secondary btn-sm" onclick="rotateRight()">R Rot</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="zoomOut()">Zoom Out</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="flipv()">V-Flip</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="grey()">Gray</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picblur(-0.5)">Sharpen</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picformat('png')">PNG</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picformat('webp')">WebP</button>
                                <button type="button" class="btn btn-secondary btn-sm" onclick="picformat('tiff')">TIFF</button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <h6>More image processing options in the <a href="https://tlnet.top/wfsdoc" target="_blank">Docs</a></h6>
                    <button type="button" class="btn btn-primary" data-bs-dismiss="modal">Close</button>
                </div>
            </div>
        </div>

        <div class="toast-container" style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); z-index: 1051;">
            <div id="copyToast" class="toast align-items-center text-white border-0" role="alert" aria-live="assertive" aria-atomic="true" data-bs-autohide="true" data-bs-delay="1800" style="background-color: rgba(0, 0, 0, 0.4);width: auto;">
                <div class="d-flex">
                    <div class="toast-body">
                        <h6 id="copyToastId" style="margin: 0;"></h6>
                    </div>
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
                        + '<td><button class="btn btn-primary m-1" onclick=\'processModalshow(' + item.Id + ',"' + aurl + item.Name + '")\';">Image Processing</button><button class="btn btn-primary m-1" onclick=\'renameshow(' + item.Id + ',"' + item.Name + '")\';">rename</button><button class="btn btn-primary" onclick=\'deletefile(this,"' + item.Name + '")\';">delete</button></td>'
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
                    alert("File Upload Sucessful:" + data.name)
                    clearFileInput();
                    document.getElementById("filepath").value = "";
                } else {
                    alert("File Upload Failed" + data.desc)
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
        if (confirm("Confirm to delete the file？")) {
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
                    alert("Failed to delete file：" + data.desc)
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
                alert("file rename failed:" + data.desc)
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

<script>
    const processModal = new bootstrap.Modal(document.getElementById('processpicModal'));

    function processModalshow(id, src) {
        processModal.show();
        document.getElementById("processpicModalLabel").innerText = "";
        document.getElementById("picsrcid").src = src;
        newImageModifier(src);
    }

    function showpic() {
        const s = imageModifier.tostring();
        if (s != "") {
            document.getElementById("picsrcid").src = s
        } else {
            alert("Not image");
        }
        document.getElementById("processpicModalLabel").innerText = s;
    }

    function zoomIn() {
        imageModifier.setZoomIn();
        showpic();
    }

    function zoomOut() {
        imageModifier.setZoomOut()
        showpic();
    }

    function rotateLeft() {
        imageModifier.setRotateLeft()
        showpic();
    }

    function rotateRight() {
        imageModifier.setRrotateRight()
        showpic();
    }

    function fliph() {
        imageModifier.setFliph()
        showpic();
    }

    function flipv() {
        imageModifier.setFlipv()
        showpic();
    }

    function picblur(b) {
        imageModifier.setBlur(b);
        showpic();
    }

    function grey() {
        imageModifier.setGrey()
        showpic();
    }

    function invert() {
        imageModifier.setInvert()
        showpic();
    }

    function picformat(fmt) {
        imageModifier.setFormat(fmt)
        showpic();
    }

    let imageModifier;

    function newImageModifier(src) {
        imageModifier = new ImageModifier(src);
    }

    class ImageModifier {
        constructor(src) {
            this.src = src;
            this.image = new Image();
            this.image.src = src;
            this.originalWidth = null;
            this.originalHeight = null;
            this.width = 0;
            this.height = 0;
            this.grey = 0;
            this.invert = 0;
            this.fliph = 0;
            this.flipv = 0;
            this.blur = 0;
            this.rotate = 0;
            this.scale = "";
            this.format = "";
            this.notimage = false;
            this.image.onload = () => {
                this.originalWidth = this.image.naturalWidth;
                this.originalHeight = this.image.naturalHeight;
                this.width = this.originalWidth;
                this.height = this.originalHeight;
            };
            this.image.onerror = () => {
                this.notimage = true;
                document.getElementById("picsrcid").src = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svgCode);
            };
        }

        setZoomIn(factor = 1.2) {
            this.width = Math.ceil(this.width + this.originalWidth * 0.1);
            this.height = Math.ceil(this.height + this.originalHeight * 0.1);
        }

        setZoomOut(factor = 0.8) {
            const w = Math.floor(this.width - this.originalWidth * 0.1);
            const h = Math.floor(this.height - this.originalHeight * 0.1);
            if (w > 0 && h > 0) {
                this.width = w;
                this.height = h;
            }
        }

        setRotateLeft() {
            this.rotate = this.rotate + 10;
        }

        setRrotateRight() {
            this.rotate = this.rotate - 10;
        }

        setGrey() {
            this.grey = this.grey == 0 ? 1 : 0;
        }

        setBlur(b) {
            this.blur = this.blur + b;
            if (this.blur < 0) {
                this.blur = 0;
            }
        }

        setFliph() {
            this.fliph = this.fliph == 0 ? 1 : 0;
        }

        setFlipv() {
            this.flipv = this.flipv == 0 ? 1 : 0;
        }

        setFormat(fmt) {
            this.format = this.format == fmt ? "" : fmt;
        }

        setInvert() {
            this.invert = this.invert == 0 ? 1 : 0;
        }

        tostring() {
            if (this.notimage) {
                return "";
            }
            let s = ""
            if (this.width != this.originalWidth || this.height != this.originalHeight) {
                s = s + "/scale/" + this.width + "x" + this.height;
            }
            if (this.grey == 1) {
                s = s + "/grey/1"
            }
            if (this.invert == 1) {
                s = s + "/invert/1"
            }
            if (this.rotate != 0) {
                s = s + "/rotate/" + this.rotate;
            }
            if (this.blur != 0) {
                s = s + "/blur/" + this.blur;
            }
            if (this.fliph == 1) {
                s = s + "/fliph/1";
            }
            if (this.flipv == 1) {
                s = s + "/flipv/1";
            }
            if (this.format != "") {
                s = s + "/format/" + this.format;
            }
            if (s != "") {
                return this.src + "?mode/0" + s;
            } else {
                return this.src
            }
        }
    }

    function copysrc() {
        copyToClipboard(document.getElementById("processpicModalLabel").innerText);
    }

    function copyToClipboard(text) {
        if (text == "") {
            return;
        }
        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).then(function () {
                console.log('Text copied to clipboard successfully!');
                copyToastShow('Text copied to clipboard');
            }).catch(function (error) {
                console.error('Failed to copy text: ', error);
                copyToastShow('Text copied to clipboard failed:' + error);
            });
        } else {
            copyToastShow('Cannot copy to clipboard, https is required');
        }
    }

    function copyToastShow(msg) {
        const toastElement = document.getElementById('copyToast');
        document.getElementById("copyToastId").innerText = msg;
        const toast = new bootstrap.Toast(toastElement);
        toast.show();
    }

</script>
</body>

</html>`

	fragmentEnText = `
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
                    <a class="nav-link" href='/init'>Account Management</a>
                    <a class="nav-link" href='/file'>File Management</a>
                    <a class="nav-link active" href='/fragment'>Fragmentation Cleanup</a>
                    <a class="nav-link" href='/monitor'>System Monitoring</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>Log In</a>
                    <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: xx-small;">
        <table class="table table-striped" style="font-size: smaller;">
            <tr>
                <th>File Name</th>
                <th>Last update time</th>
                <th>Size (Bytes)</th>
                <th>Fragment size(Bytes)</th>
                <th>File status</th>
                <th>File Operations</th>
            </tr>
            <tbody id="fileTableBody">
                {{range $k,$v := .}}
                <tr>
                    <td>{{$v.Name}}</td>
                    <td>{{$v.Time}}</td>
                    <td>{{$v.FileSize}}</td>
                    <td>{{$v.FragmentSize}}</td>
                    {{if eq $v.Status 1 }}
                    <td>readonly</td>
                    {{else}}
                    <td>Read and Write(Forbidden)</td>
                    {{end}}
                    {{if and (eq $v.Status 1) (gt $v.FragmentSize 0)}}
                    <td><button id="{{$v.Name}}" class="btn btn-primary" onclick="fragment('{{$v.Name}}')">defragment</button>
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
            if (confirm("[" + n + "]Are you sure you want to defragment? \nIt is recommended that defragmentation should be performed in a state where WFS service operations are relatively low to reduce the impact on front-end service quality ")) {
                document.getElementById(n).innerText = "defragmenting..."
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
                    document.getElementById(n).innerText = "defragment"
                    document.getElementById(n).disabled = false;
                    return response.json();
                }).then(data => {
                    if (data.status) {
                        window.location.reload(true);
                    } else {
                        alert("defragment failed")
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

	monitorEnText = `
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
                    <a class="nav-link" href='/init'>Account Management</a>
                    <a class="nav-link" href='/file'>File Management</a>
                    <a class="nav-link" href='/fragment'>Fragmentation Cleanup</a>
                    <a class="nav-link active" href='/monitor'>System Monitoring</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>Log In</a>
                    <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="container mt-1 card" style="font-size: small;">
        <div class="container mt-1" style="font-size: small;">
            <h3>performance data monitoring</h3>
            <div class="input-group">
                <span class="input-group-text">time interval(unit:second)</span>
                <input id="stime" placeholder="time interval" value="3" />
                <button class="btn btn-primary" onclick="monitorLoad();">start</button>&nbsp;
                <button class="btn btn-primary" onclick="stop();">stop</button>&nbsp;
                <button class="btn btn-primary" onclick="clearData();">data clear</button>
            </div>
        </div>

        <table class="table table-striped " style="font-size: smaller;">
            <tr>
                <th></th>
                <th>memory usage(MB)</th>
                <th>Memory allocated(MB)</th>
                <th>Memory release times</th>
                <th>Coroutine number</th>
                <th>CPU number</th>
                <th>Disk Free Space(GB)</th>
                <th>Memory usage</th>
                <th>CPU usage</th>
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

	aboutEnText = `<html>

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
                        <a class="nav-link" href='/init'>Account Management</a>
                        <a class="nav-link" href='/file'>File Management</a>
                        <a class="nav-link" href='/fragment'>Fragmentation Cleanup</a>
                        <a class="nav-link" href='/monitor'>System Monitoring</a>
                    </div>
                    <div class="navbar-nav ms-auto">
                        <a class="nav-link" href='/login'>Log In</a>
                        <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                    </div>
                </div>
            </div>
        </nav>
        <div class="container my-1 card">
            <h3>WFS Introduction</h3>
            <h6>Copyright, Author <a href="mailto:donnie4w@gmail.com" class="text-reset">donnie4w</a>, version {{.}} </h6>
            <hr>
            <div class="my-2">
                <p><strong>WFS File Storage System</strong>，primarily used to handle the storage of a large number of small files</p>
                <p>WFS Official Website:<a href="https://tlnet.top/wfs" target="_blank">https://tlnet.top/wfs</a></p>
                <p><a href="https://tlnet.top/wfsdoc" target="_blank">WFS User Documentation</a></p>
            </div>
            <hr>
            <h5 class="mt-2">wfs Related Programs</h5>
			<div class="row button-links">
				<div class="col-sm-6 col-md-4 col-lg-3 mb-3">
					<a href="https://github.com/donnie4w/wfs" target="_blank" class="btn btn-primary">Wfs Source Code</a>
				</div>
				<div class="col-sm-6 col-md-4 col-lg-2 mb-3">
					<a href="https://github.com/donnie4w/wfs-goclient" target="_blank" class="btn btn-secondary">Go Clients</a>
				</div>
				<div class="col-sm-6 col-md-4 col-lg-2 mb-3">
					<a href="https://github.com/donnie4w/wfs-rsclient" target="_blank" class="btn btn-secondary">Rust Clients</a>
				</div>
				<div class="col-sm-6 col-md-4 col-lg-2 mb-3">
					<a href="https://github.com/donnie4w/wfs-jclient" target="_blank" class="btn btn-secondary">Java Clients</a>
				</div>
				<div class="col-sm-6 col-md-4 col-lg-2 mb-3">
					<a href="https://github.com/donnie4w/wfs-pyclient" target="_blank" class="btn btn-secondary">Python
						Clients</a>
				</div>
			</div>
    
            <p class="text-muted mt-3">
                <strong>Email:</strong>
                <a href="mailto:donnie4w@gmail.com" class="text-reset">donnie4w@gmail.com</a>
            </p>
        </div>
    
    </html>`
)
