/*************************************************************************
  * File Name :  node.js
  * Author  :      unasm
  * Mail :         unasm@sina.cn
  * Last_Modified: 2015-05-25 09:56:26
 ************************************************************************/

console.log("Server started");
var connects = new Array();
var WebSocketServer = require('ws').Server
, wss = new WebSocketServer({port: 8010});
console.log("complete");
wss.on('connection', function(ws) {
	connects.push(ws);
	ws.on('message', function(msg) {
		ws.send(msg);
		for(var i = 0,len = connects.length; i < len;i++){
			if(connects[i] !== ws){
				connects[i].send(msg);
			}
		}
	});
	console.log("connection");
});
