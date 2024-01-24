# Dev 

- Goals:
    - Complete the webserver
        - What is manditory and not - if not add to dev/{feature}
    - Move all wild ideas in the dev/{feature} structure
    - Any more wild ideas go in he appropriate directory
    - What do I need to do?
        - args 
        - make cli app part modular - so I do not have to program it - to the app without breaking anything 
        - handler - serverTypes WebServer http, https !! DO NOT ADD ANYTHING ELSE
        - TLS
            - openssl command
            - import a cert
    - End up at the winchester have the next update ideas in order of:
        - console && cmds Ninja Shell be mergable into the cli code for the next major update.
        - ProxyServer type
        - MITMListener type
        - RogueDNS/LDAP     


Flow
- main -> handleArgs 
- handleArgs.go marshals data to struct (to initStruc)
    - SetDefaultConfig() to force no nil values
    - InitSeverFromArgs (instead from InitServerStruct)
    - metahandler.SelectAction(s *Server runFlag) 
- metahandler (to do the heavy lifting )
    - SelectAction
    - Actions...
    - metahandler  and create/start/stop/config passing ServerType and some flag to (console would fit in here: either run create server or goto Console)
    - NOT SURE ABOUT CONFIG
- web ( because web is the priority)
- http
- either http OR tls and https

Concerns:
- Going overboard with args
- Logging




Todo
1. main.go : args - make new main, compare struct in metahandler, test variables in main to create commands:
    - http
    - https   
2. metahandler : args -> marshalToStruct -> Selection -> create/start/stop/config -> Graceful Exit
3. Web -> http / https - if I start by categorising at this level it make the above less cluttered and forces modularity of the sub categorises
4. HTTP
5. HTTPS <-> TLS


Comments cleaning
```go

//
//
// (NAME OF CONCEPT THAT MANAGES) -> server1,server2,...
// Seperation of the methods as I am double server and IDdatabase
// X-server: web-server.go, proxy-server.go
// CURRENT IDEA Database needs to be:
// - part of larger struct that: map[string](pointer) points to Server structs, ID database etc 
// - initialisation of array to make ID database - ID need negative space for stopped servers
// 

// Are negative ID is a good way of managing this why not flags
//
       
// 
// Creation to termination 
// Memory Arenas

// 
// IDs
// Memory Arenas and how.. 
// 
```