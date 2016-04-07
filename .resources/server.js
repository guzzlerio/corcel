#! /usr/bin/env node

const http = require('http');
const url = require("url");
const hostname = 'localhost';
const port = 1337;

var counter = 0;

var codes = [200,400,500];

http.createServer((req, res) => {
      res.writeHead(codes[counter++ % 3], { 'Content-Type': 'text/plain' });
      require('crypto').randomBytes(1 , function(err, buffer) {
            var token = buffer.toString('hex');
            res.end(token);
      });
      //console.log(url.parse(req.url).pathname, 200)
}).listen(port, hostname, () => {
      console.log(`Server running at http://${hostname}:${port}/`);
});
