
-- 用户表
CREATE TABLE `users`(
    `id` int(10) unsigned PRIMARY KEY auto_increment,
    `login_name` varchar(64) NOT NULL UNIQUE KEY,
    `pwd` text NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- 视频表
CREATE TABLE `videoinfo`(
    `id` varchar(64) PRIMARY KEY,
    `authorid` int(10),
    `name` text,
    `display_ctime` text,
    `create_time` datetime DEFAULT current_timestamp
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- 评论表
CREATE TABLE `comments`(
    `id` varchar(64) PRIMARY KEY,
    `videoid` varchar(64),
    `authorid` int(10),
    `content` text,
    `time` datetime DEFAULT current_timestamp
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- session表
CREATE TABLE `sessions`(
    `sessionid` varchar(256) PRIMARY KEY,
    `TTL` text,
    `loginname` text
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
