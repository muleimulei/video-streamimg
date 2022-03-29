$(document).ready(function(){

    currentVideo = null
    listedVideos = null
    uid = 0

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

    initPage(null)


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
                    // console.log("empty videos")
                    return
                }
                listedVideos = obj['videos'];
                obj['videos'].forEach(function(item, index) {
                    var ele = htmlVideoListElement(item['id'], item['name'], item['display_ctime']);
                    $("#items").append(ele);
                });
                callback();
            });
        });
    }

    function htmlVideoListElement(vid, name, ctime) {
        var ele = $('<a/>', {
            href: '#'
        });
        ele.append(
            $('<video/>', {
                width:'320',
                height:'240',
                poster:'/statics/img/preloader.png',//20200821
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

    // Comments operations
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
            // console.log(res, textStatus, xhr )
            if (xhr.status == 200) {
                callback(res, null);
                return;
            }
        }
    })
}


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
            // console.log(res, textStatus, xhr )
            if (xhr.status == 200) {
                callback(res, null)
                return
            }
        }
    })
}
})