# telebot-worker

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

You command have to read the json data from STDIN, and write reply message to STDOUT.

## php hello world example

```php
<?php
// this file is named "handler.php"

$data = json_decode(file_get_contents("php://input"), true);
echo "Hello, " + $data["from"]["first_name"] + " " + $data["from"]["last_name"];
```

```sh
./telebot-worker -worker 5 -token mytoken.txt php handler.php
```

## Special thanks to

* [Telebot](https://github.com/tucnak/telebot)
* [Telegram](https://telegram.org)

## License
Any version of MIT, GPL or LGPL
