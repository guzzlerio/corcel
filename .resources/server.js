#! /usr/bin/env node

const http = require('http');
const url = require("url");
const hostname = 'localhost';
const port = 1337;

http.createServer((req, res) => {
      res.writeHead(200, { 'Content-Type': 'text/plain' });
      res.end('Hello World\n');
      //console.log(url.parse(req.url).pathname, 200)
}).listen(port, hostname, () => {
      console.log(`Server running at http://${hostname}:${port}/`);
});
