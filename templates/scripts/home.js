$(document).ready(function(){
    uname = ''
    session = ''
    DEFAULT_COOKIE_EXPIRE_TIME = 5


    $('#signinhref').on('click', function() {
        $('#signin').show()
        $('#signup').hide()
    })

    $('#signuphref').on('click', function() {
        $('#signup').show()
        $('#signin').hide()
    })

    function popupErrorMsg(e, err){
        e.text(err)
        setTimeout(function(){e.text("")}, 2000)
    }

    function setCookie(cname, cvalue, exptime){
        var d = new Date();
        d.setTime(d.getTime() + DEFAULT_COOKIE_EXPIRE_TIME * 60 * 1000)
        var expires = "expires=" + d.toUTCString()
        document.cookie = cname + "=" + cvalue+";" + expires+";path=/"
    }




    //用户注册
    function registerUser(callback) {
        var username = $('#username1').val()
        var pwd = $('#Password1').val()

        if(username == '' || pwd == ''){
            callback('', '用户名和密码不能为空')
            return
        }

        var apiUrl = window.location.hostname + ":8080/api"

        var reqBody = {
            "user_name" : username,
            "pwd": pwd,
        }

        var da = {
            'url': 'http://' + window.location.hostname + ':8000/user',
            'method': 'POST',
            'req_body': JSON.stringify(reqBody),
        }


        $.ajax({
            url: 'http://' + apiUrl,
            type: 'POST',
            data: JSON.stringify(da),
            success: function (data, statusText, xhr) {
                console.log(data, statusText, xhr)
                if(xhr.status == 200) {
                    uname = username
                    callback(data, null)
                }else{
                    callback(null, 'Error of register')
                }
            },
            error: function(err) {
                console.log(err)
                g = JSON.parse(err.responseText)
                callback(null, g.error)
            }
        })
    }


    //注册
    $('#signupbtn').on('click', function(e) {
        $('#signupbtn').text('loading ...')
        e.preventDefault()
        registerUser(function(res, err) {
            $('#signupbtn').text('注册')
            if(err != null){
                popupErrorMsg( $('#bg-warning1'), err)
                return
            }

            var r = JSON.parse(res)
            setCookie('session', r['session_id'], DEFAULT_COOKIE_EXPIRE_TIME)
            setCookie('username', uname, DEFAULT_COOKIE_EXPIRE_TIME)

            $('#signupsubmit').submit()
        })
    })


    //用户登录
    function loginUser(callback) {
        var username = $('#username2').val()
        var pwd = $('#Password2').val()

        if(username == '' || pwd == ''){
            callback('', '用户名和密码不能为空')
            return
        }

        var apiUrl = window.location.hostname + ":8080/api"

        var reqBody = {
            "user_name" : username,
            "pwd": pwd,
        }

        var da = {
            'url': 'http://' + window.location.hostname + ':8000/user/' + username,
            'method': 'POST',
            'req_body': JSON.stringify(reqBody),
        }

        $.ajax({
            url: 'http://' + apiUrl,
            type: 'POST',
            data: JSON.stringify(da),
            success: function (data, statusText, xhr) {
                console.log(data, statusText, xhr)
                if(xhr.status == 200) {
                    uname = username
                    callback(data, null)
                }else{
                    callback(null, 'Error of register')
                }
            },
            error: function(err) {
                console.log(err)
                g = JSON.parse(err.responseText)
                callback(null, g.error)
            }
        })

    }

    //登录
    $('#signinbtn').on('click', function(e) {
        $('#siginbtn').text('loading ...')
        e.preventDefault()
        loginUser(function(res, err) {
            $('#signupbtn').text('登录')
            if(err != null){
                popupErrorMsg( $('#bg-warning2'), err)
                return
            }

            var r = JSON.parse(res)
            setCookie('session', r['session_id'], DEFAULT_COOKIE_EXPIRE_TIME)
            setCookie('username', uname, DEFAULT_COOKIE_EXPIRE_TIME)

            $('#signinsubmit').submit()
        })
    })
})

