var request = require('request');

var body = {
    "to":[
        {
            "id": "b4b19d9a4337cad8631299d6519aa17a1fe58d145ad0235dadbc5c8ee2e1a7bf",
            "user": "everhusk"
        }
    ],
    "body": "testing",
    "type": "text"
};

var options = {
  method: 'POST',
  body: body,
  json: true,
  url: 'https://66jezutq1k.execute-api.us-east-1.amazonaws.com/production/v1/update/kik',
  headers: {
    'Content-Type': 'application/json',
  },
}

request(options, function (err, res, body) {
  if (err) {
    console.log('err',err)
    return
  }
  //console.log('res',res);
  console.log('body',body)
});