const dgram = require('dgram');

const SERVER_ADDR = '2001:412:abcd:2:0013:A200:4147:8C2B';
const SERVER_PORT = 8000;

const client = dgram.createSocket('udp6');

var c = 0;

function sendData() {
    const message = Buffer.from('HELLO ' + c);
    client.send(message, SERVER_PORT, SERVER_ADDR);
    c++;
}

setInterval(sendData, 1000);
