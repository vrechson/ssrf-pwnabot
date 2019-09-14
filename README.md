# ssrf pwnabot  
discover secret services and pwn this bot with ssrf  

env example:  

WEBHOOK_SERVICE_URL=https://url-api.herokuapp.com  
TELEGRAM_API_TOKEN=your_token  
PORT=8083 // do not define if you're using heroku  
SERVICE_PORT=8080 // secret service port  
DEBUG=true  
USE_WEBHOOK=true // true for heroku but false if you want polling or run locally tests  