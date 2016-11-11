const http = require('http');
const url = require("url");
const hostname = 'localhost';
const port = 1337;

http.createServer((req, res) => {
    res.writeHead(200, {
        'Content-Type': 'text/plain'
    });
    res.end();

}).listen(port, hostname, () => {
    console.log(`Server running at http://${hostname}:${port}/`);

});
