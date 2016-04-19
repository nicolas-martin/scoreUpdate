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
var request = require('sync-request');
var NHL = require('NHL');

// Find all teams that match the User message
var getTeamName = function(message,type){
    // Iterate over each word
    for(var i=message.length-1;i>0;i--){
        // Check NHL teams
        var matches=[];var match;
        for(var j=0;j<NHL.teams.length;j++){
            // Add team to suggestions
            teamName = NHL.teams[j].body.toUpperCase().split(" ");
            // Iterate over each word on the team to see if it matches
            for(var k=teamName.length-1;k>=0;k--){
                if(teamName[k] == message[i]){
                    match = {
                        name: NHL.teams[j].body,
                        id: NHL.teams[j].id
                    };
                    matches.push({
                        "type": "text",
                        "body": type+" "+NHL.teams[j].body
                    });
                }
            }
        }
        // Return results
        if(matches.length == 1){
            return match;
        }
    }
    return matches;
}

// Main Run
exports.handler = function(event, context){
    // Initilizations
    var from,body,chatId;
    var messages = event.body.messages;
    var message = "";
    var responses = [];
    // Multiple messages at once might come in
    for(var i=0;i<messages.length;i++){
        // Check if welcome message
        if(messages[i].type == "start-chatting"){
            var body = "Hey, "+messages[i].from+"! Welcome to Sports Update. To subscribe to a team type 'Subscribe <Team Name>', to unsubscribe type 'Unsubscribe <Team Name>.";
            var suggestions = false;
            // Register the user
            var options = { 
                headers: {'Content-Type': 'application/json'},
                json: {
                    "Username": messages[i].from,
                    "chatId": message[i].chatId,
                    "Platform":"Kik"
                }
            };
            var res = request('POST','https://sportsbot-1255.appspot.com/User',options);
            if(res.statusCode >= 300){
                console.log('ERR-000',res.body.toString()+". Failed to setup user "+messages[i].from);
            }
        }
        // Check if text message
        else if(messages[i].type == 'text'){
            // Standardize message characterization (upper case)
            message = messages[i].body.toUpperCase();
            // Attempt to subscribe a user from a team
            if(message.indexOf("SUBSCRIBE") > -1){
                // Check if it's an unsubscribe
                var action = (message.indexOf("UNSUBSCRIBE") > -1) ? "Unsubscribe" : "Subscribe";
                // Find matching team names
                match = getTeamName(message.split(" "),action);
                // Multiple or no matching teams found
                if(match instanceof Array){
                    if(match.length > 0){
                        var body = "Please select a team";
                        var suggestions = teamName;
                    }else{
                        var body = "No team found, please try again.";
                        var suggestions = false;
                    }
                }
                // Team name found. Try to subscribe/unsubscribe User.
                else{
                    var suggestions = false;
                    var options = { 
                        headers: {'Content-Type': 'application/json'},
                        json: {"Username":messages[i].from,"TeamID":match.id}
                    };
                    var res = request('POST','https://sportsbot-1255.appspot.com/'+action,options);
                    if(res.statusCode >= 300){
                        console.log("ERR-001",res.body.toString());
                        var body = action+" failed. Please try again shortly."
                    }
                    else{
                        console.log('body',res.body.toString());
                        var body = "Successfully "+action+"d to the "+match.name+".";
                    }
                }
            }
            else{
                // Default response
                var body = "Subscribe to a team with 'Subscribe team name' and unsubscribe with 'Unsubscribe team name'";
                var suggestions = false;
            }
        }
        // Invalid message format
        else{
            break;
        }
        // Format the response for the Kik bot
        responses.push({
            'body': body, 
            'to': messages[i].from, 
            'type': 'text', 
            'chatId': messages[i].chatId
        });
        if(suggestions){
            responses.keyboards = [{
                "type": "suggested",
                "responses": suggestions
            }];
        }
    }

    // Return all responses
    var options = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Basic ' + new Buffer('sportsupdate' + ':' + '09439681-97b4-4d96-b51a-a03549a76ede').toString('base64')
        },
        json:{messages:responses} 
    };
    var res = request('POST','https://api.kik.com/v1/message',options);
    if(res.statusCode >= 300){
        console.log("ERR-002",res.body.toString());
        context.fail("Message Response Failed");
    }
    else{
        console.log('body',res.body.toString());
        context.succeed("Success");
    }
};