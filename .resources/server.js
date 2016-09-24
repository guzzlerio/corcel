#! /usr/bin/env node

const http = require('http');
const url = require("url");
const crypto = require('crypto');

const hostname = 'localhost';
const port = 1337;

var counter = 0;

var codes = [200, 400, 500];

var responders = {
    "slow": (req, res) => {
        var min = 1000;
        var max = min * 10;
        setTimeout(() => {
            res.writeHead(200, {
                'Content-Type': 'text/plain'
            });
            res.end();
        }, randomInt(min, max));
    },
    "fail": (req, res) => {
        res.writeHead(500, {
            'Content-Type': 'text/plain'
        });
        res.end();
    },
    "ok": (req, res) => {
        res.writeHead(200, {
            'Content-Type': 'text/plain'
        });
        res.end();
    },
    "big": (req, res) => {
        var min = 1024 * 1024;
        var max = min * 20;
        var size = randomInt(min, max);
        crypto.randomBytes(size, function(err, buffer) {
            var token = buffer.toString('hex');
            res.writeHead(200, {
                'Content-Type': 'text/plain'
            });
            res.end(token);
        });
    }
};

function randomInt(low, high) {
    return Math.floor(Math.random() * (high - low) + low);
}

function selectResponder() {
    var funcs = Object.keys(responders).map((key) => {
        return responders[key];
    });
    var min = 0;
    var max = funcs.length;
    var selection = randomInt(min, max);
    return funcs[selection];
}

var responder = selectResponder();

http.createServer((req, res) => {
    var responder = selectResponder();
    responder(req, res);
}).listen(port, hostname, () => {
    console.log(`Server running at http://${hostname}:${port}/`);
});
