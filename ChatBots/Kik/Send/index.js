// ************************************************************************
//                     KIK CHATBOT Function
//
//                      INPUTS: Message
//             OUTPUTS: Posts Message to Kik User (SUCCESS || FAIL)
//
//                            Overview:
//
//      1) Parse Incoming Message
//      2) Respond with clever message
// ************************************************************************
var AWS = require('aws-sdk');
AWS.config.update({accessKeyId: '', secretAccessKey: ''});
var request = require('request');
var async = require("async");

// Main Run
exports.handler = function(event, context){
    // Initilizations
    var to = event.body.to;
    var body = event.body.body;
    var type = event.body.type;
    var queue = [];
    var responses = [];
    // Format responses
    for(var i=0;i<to.length;i++){
        responses.push({
            'body': body,
            'to': to[i].user,
            'type': 'text',
            'chatId': to[i].id
        });
        // If 25 message batch is filled
        if(responses.length == 25){
            queue.push({responses:responses}); // Push batch to queue
            responses = []; // Reset responses array
        }
    }
    queue.push({responses:responses}); // Push remaining responses to queue
    // Setup messaging queue
    var q = async.queue(function (task, callback) {
        console.log('Sending Message:',task.responses);
        // Return all responses
        var options = {
            method:'POST',
            uri:'https://api.kik.com/v1/message',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Basic ' + new Buffer('sportsupdate' + ':' + '09439681-97b4-4d96-b51a-a03549a76ede').toString('base64')
            },
            json:{messages:task.responses}
        };
        request(options,callback);
    }, 2);
    // Setup complete callback
    q.drain = function() {
        context.succeed("Success");
    }
    console.log('Queuing:',queue);
    // Add responses to the queue
    q.push(queue,function(err){
        if(err){
            console.log("ERR-002",res);
        }
    });
};