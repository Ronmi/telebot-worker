/*
telebot-worker is a simple telegram bot wrapper focus on text messages.

This program is mean to provide aids to people who writing php or any other language which does not suitable to make telegram bot.

Build step for lazy people

You have to authorize your bot and put the token in a text file before running this program. Then:

    export GOPATH=/path/to/your/workspace
    go get github.com/Ronmi/telebot-worker
    go build github.com/Ronmi/telebot-worker
    bin/telebot-worker -worker 5 -token /path/to/token/file cmd_to_run args...

Message passing

This program will run command in sub-process, and feed received telegram messge in json format through STDIN.
Your command have to read the json data from STDIN, and write reply message (json format) to STDOUT.

Commandline parameters

There are two parameters in telebot-worker.

    worker:
        Number of worker thread.
    token:
        path to token file.

The worker threads defaults to 5.

Example

Here shows an example app which sends a photo when you text to it.
    <?php
    // this file is named "handler.php"

    $data = json_decode(file_get_contents("php://stdin"), true);
    $ret = array(
        "content" => "/path/to/hello/world.jpg"
        "type" => "photo"
    );
    echo json_encode($ret)

Then run with
    bin/telebot-worker -worker 5 -token mytoken.txt php handler.php
*/
package main
