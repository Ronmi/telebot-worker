# telebot-worker

[![Build Status](https://travis-ci.org/Ronmi/telebot-worker.svg?branch=master)](https://travis-ci.org/Ronmi/telebot-worker)

A simple telegram bot wrapper focus on text messages.

This program is mean to provide aids to people who writing php or any other language which does not suitable to make telegram bot.

## Synopsis

```sh
go get
go build
./telebot-worker -worker 5 -token /path/to/token/file cmd_to_run args...
```

## Message passing

This program will run command in sub-process, and feed received telegram messge in json format through STDIN.

You command have to read the json data from STDIN, and write reply message (json format) to STDOUT.

## Return in json format

Returned data have to be an array of objects, which contain these fields:

 - type: Message type. String of "text", "doc", "audio", "video", "photo"
 - user: Recipent. See [User](https://core.telegram.org/bots/api/#user) and [GroupChat](https://core.telegram.org/bots/api/#groupchat)
 - content: Message body for text message, or filename for other type of message.

You can omit `user` field for replying to original user, or omit type for text message.

## php hello world example

```php
<?php
// this file is named "handler.php"

$data = json_decode(file_get_contents("php://stdin"), true);
$ret = array(array(
    "content" => "/path/to/hello/world.jpg"
    "type" => "photo"
));
echo json_encode($ret)
```

```sh
./telebot-worker -worker 5 -token mytoken.txt php handler.php
```

## Special thanks to

* [Telebot](https://github.com/tucnak/telebot)
* [Telegram](https://telegram.org)

## License
Any version of MIT, GPL or LGPL
