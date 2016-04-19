var request = require('request');

var formData = {
    "webhook": "https://66jezutq1k.execute-api.us-east-1.amazonaws.com/production/v1/kik",
    "features": {
       "manuallySendReadReceipts": false,
       "receiveReadReceipts": false,
       "receiveDeliveryReceipts": false,
       "receiveIsTyping": false
    }
}

var options = {
  method: 'POST',
  form: JSON.stringify(formData),
  url: 'https://api.kik.com/v1/message',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Basic ' + new Buffer('sportsupdate' + ':' + '09439681-97b4-4d96-b51a-a03549a76ede').toString('base64')
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