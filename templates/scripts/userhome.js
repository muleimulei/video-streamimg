$(document).ready(function(){
    currentVideo = null
    listedVideos = null
    uid = 0 //用户id

    //错误信息
    function popupErrMsg(msg){
        $('#errorbar'),text(msg);
        $('#errorbar').show()

        setTimeout(() => {
            $('#errorbar').hide()
        }, 2000);
    }

    //正确信息
    function popNotificationMsg(msg){
        $('#snackbar').text(msg);
        $('#snackbar').show()
        setTimeout(() => {
            $('#snackbar').hide()
        }, 2000);
    }

    // 获取cookie
    function getCookie(cname) {
        var name = cname + "="
        var ca = document.cookie.split(';')
        for(var i = 0; i < ca.length; i++) {
            var c = ca[i]
            while(c.charAt(0) == ' ') {
                c = c.substring(1);
            }
            if(c.indexOf(name) == 0) {
                return c.substring(name.length, c.length)
            }
        }
        return ""
    }
    //获取用户id
    function getUserId(callback) {
        var apiUrl = window.location.hostname + ":8080/api"
        var dat = {
            'url': 'http://' + window.location.hostname + ':8000/user/' + getCookie('username'),
            'method': 'GET'
        };

        $.ajax({
            url: 'http://' + apiUrl,
            type: 'post',
            data: JSON.stringify(dat),
            headers: {'X-Session-Id': getCookie('session')},
            error: function(xhr){
                // console.log(xhr)
                g = JSON.parse(xhr.responseText)
                callback(null, g.error)
            },
            success: function(res, textStatus, xhr) {
                // console.log(res, textStatus, xhr )
                if (xhr.status == 200) {
                    callback(res, null);
                    return;
                }
            }
        })
    }

    function listAllVideos(callback) {
        var apiUrl = window.location.hostname + ":8080/api"
        var dat = {
            'url': 'http://' + window.location.hostname + ':8000/user/' + getCookie('username') + '/videos',
            'method': 'GET',
        }
    
        $.ajax({
            url: 'http://' + apiUrl,
            type: 'post',
            data: JSON.stringify(dat),
            headers: {'X-Session-Id': getCookie('session')},
            error: function(xhr){
                // console.log(xhr)
                g = JSON.parse(xhr.responseText)
                callback(null, g.error)
            },
            success: function(res, textStatus, xhr) {
                // console.log(res, textStatus, xhr )
                if (xhr.status == 200) {
                    callback(res, null);
                    return;
                }
            }
        })
    }

    //选择播放哪个视频
    function selectVideo(vid) {
        var url = 'http://' + window.location.hostname + ':9000/videos/'+vid
        $('#curr-video').attr('src', url)
        $('#curr-video-name').text(currentVideo['name'])
        $('#curr-video-ctime').text('Uploaded at: '+ currentVideo['display_ctime'])

        refreshComment(vid)
    }

    initPage(function() {
        if(listedVideos !== null) {
            currentVideo = listedVideos[0];
            //选择第一个视频
            selectVideo(listedVideos[0]['id'])
        }

        $('.video-item').click(function(){
            var self = this.id
            listedVideos.forEach(function(item, index) {
                if(item['id'] === self) {
                    currentVideo = item
                    return
                }
            })
            selectVideo(self)
        })

        $('#submit-comment').click(function(){
            if(currentVideo == null) return
            var content = $('#comments-input').val()

            postComment(currentVideo['id'], content, function(res, err){
                if(err != null) {
                    popupErrMsg("Error when try to post a comment")
                    return
                }

                if(res === 'ok') {
                    popNotificationMsg("New comment posted")
                    $('#comments-input').val('');
                    refreshComment(currentVideo['id'])
                }
            })
        })

        $('#uploadfile').click(function(){
            console.log('File upload')
            var filename = $('#InputFilename').val()
            var file = $('#InputFile').val()

            if(filename == '' || file == '')
            {
                $('#bg-warning3').text('filename and file can not be empty')
                setTimeout(function(){ $('#bg-warning3').text('')  }, 2000)
            }

            createVideo(filename, function(res, err){
                if(err != null){
                    $('#bg-warning3').text('Error when try to upload video')
                    setTimeout(function(){ $('#bg-warning3').text('')  }, 2000)
                    return
                }
                var obj = JSON.parse(res)

                $.ajax({
                    url : 'http://' + window.location.hostname + ':8080/upload/' + obj['id'],
                    type : 'POST',
                    data : new FormData($('#uploadform')[0]),
                    //headers: {'Access-Control-Allow-Origin': 'http://127.0.0.1:9000'},
                    crossDomain: true,
                    processData: false,  // tell jQuery not to process the data
                    contentType: false,  // tell jQuery not to set contentType
                    success : function(data) {
                        console.log(data);
                        $('#uploadvideomodal').hide();
                        location.reload();
                        //window.alert("hoa");
                    },
                    complete: function(xhr, textStatus) {
                        if (xhr.status === 204) {
                            window.alert("finish")
                            return;
                        }
                        if (xhr.status === 400) {
                            $("#uploadvideomodal").hide();
                            popupErrorMsg('file is too big');
                            return;
                        }
                    }
                });
            })
        })
    })


    function initPage(callback) {
        
        getUserId(function(res, err) {
            if (err != null) {
                window.alert("Encountered error when loading user id");
                return;
            }
            var obj = JSON.parse(res);
            uid = obj['id'];

            listAllVideos(function(res, err) {
                if (err != null) {
                    console.log(err)
                    //window.alert('encounter an error, pls check your username or pwd');
                    //popupErrorMsg('encounter an error, pls check your username or pwd');
                    return;
                }
                var obj = JSON.parse(res);
                if(obj.videos == null) {
                    console.log("empty videos")
                    return
                }
                listedVideos = obj['videos'];
                
                obj['videos'].forEach(function(item, index) {
                    var ele = htmlVideoListElement(item['id'], item['name'], item['display_ctime']);
                    $("#items").append(ele);
                });
                callback()
            });
        });
    }

    function htmlVideoListElement(vid, name, ctime) {
        var ele = $('<a/>', {
            href: '#'
        });
        ele.append(
            $('<video/>', {
                width:'400',
                height:'240',
                poster:'/static/img/preloader.png',//20200821
                controls: true
                //href: '#'
            })
        );
        ele.append(
            $('<div/>', {
                text: name
            })
        );
        ele.append(
            $('<div/>', {
                text: ctime
            })
        );
    
    
        var res = $('<div/>', {
            id: vid,
            class: 'video-item'
        }).append(ele);
    
        res.append(
            $('<button/>', {
                id: 'del-' + vid,
                type: 'button',
                class: 'del-video-button',
                text: 'Delete'
            })
        );
    
        res.append(
            $('<hr>', {
                size: '2'
            }).css('border-color', 'grey')
        );
    
        return res;
    }

    function htmlCommentListElement(cid, author, content) {
        var ele = $('<div/>', {
            id: cid
        });
    
        ele.append(
            $('<div/>', {
                class: 'comment-author',
                text: author + ' says:'
            })
        );
        ele.append(
            $('<div/>', {
                class: 'comment',
                text: content
            })
        );
    
        ele.append('<hr style="height: 1px; border:none; color:#EDE3E1;background-color:#EDE3E1">');
    
        return ele;
    }


    //刷新评论
    function refreshComment(vid) {
        listAllComments(vid, function (res, err){
            if(err != null) {
                popupErrMsg('Error when loading comments')
                return
            }
            
            var obj = JSON.parse(res)
            console.log(obj)
            $('#comments-history').empty()

            if(obj['comments'] === null) {
                $('#comments-total').text('0 Comments')
                return
            } else {
                $('#comments-total').text(obj['comments'].length + ' Comments')
            }

            obj['comments'].forEach(function(item, index){
                var ele = htmlCommentListElement(item['id'], item['author'], item['content'])
                $('#comments-history').append(ele)
            })
        })
    }

    // 提交评论
    function postComment(vid, content, callback) {
        var apiUrl = window.location.hostname + ":8080/api"

        var reqBody = {
            'author_id': uid,
            'content': content
        }

        var dat = {
            'url': 'http://' + window.location.hostname + ':8000/videos/' + vid + '/comments',
            'method': 'POST',
            'req_body': JSON.stringify(reqBody)
        };

        $.ajax({
            url  : 'http://' + apiUrl,
            type : 'post',
            data : JSON.stringify(dat),
            headers: {'X-Session-Id': getCookie('session')},
            error: function(xhr){
                // console.log(xhr)
                g = JSON.parse(xhr.responseText)
                callback(null, g.error)
            },
            success: function(res, textStatus, xhr) {
                callback(res, null);
            }
        })
    }

    // 查找该音频所有评论
    function listAllComments(vid, callback) {
        var apiUrl = window.location.hostname + ":8080/api"

        var dat = {
            'url': 'http://' + window.location.hostname + ':8000/videos/' + vid + '/comments',
            'method': 'GET',
        };

        $.ajax({
            url  : 'http://' + apiUrl,
            type : 'post',
            data : JSON.stringify(dat),
            headers: {'X-Session-Id': getCookie('session')},
            error: function(xhr){
                // console.log(xhr)
                g = JSON.parse(xhr.responseText)
                callback(null, g.error)
            },
            success: function(res, textStatus, xhr) {
                callback(res, null)
            }
        })
    }

   //上传音频


   function createVideo(vname, callback){
    var apiUrl = window.location.hostname + ":8080/api"
        var reqBody = {
            'author_id' : uid,
                'name': vname,
        }

       var dat = {
            'url': 'http://' + window.location.hostname + ':8000/user/' + getCookie('username') + '/videos',
            'method': 'POST',
            'req_body': JSON.stringify(reqBody)
        }
        $.ajax({
            url: 'http://' + apiUrl,
            type: 'post',
            data: JSON.stringify(dat),
            headers: {'X-Session-Id': getCookie('session')},
            error: function(xhr){
                // console.log(xhr)
                g = JSON.parse(xhr.responseText)
                callback(null, g.error)
            },
            success: function(res, textStatus, xhr) {
                callback(res, null)
            }
        })

   }
})