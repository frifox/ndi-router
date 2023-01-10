# Go NDI Router

NDI router whose input/output matrix can be controlled via rest or websocket. Bitfocus Companion + Stream Deck is a good use-case for this.

## ndi-router.conf
```
ApiAddr = :00
SqliteDsn = "file:ndi-router.db?_mutex=full"

[outputs]
; output_id = Channel Name
main-screen = Front Screen
rear-screen = Rear Screen
kids = Kids Room TV
kitchen = Kitchen TV
livestream = Live Stream

[inputs]
; input_id = HOSTNAME (Channel Name)
cam1 = ENCODER (Cam1)
cam2 = ENCODER (Cam2)
cam3 = ENCODER (Cam3)
pc1 = ENCODER (PC1)
pc2 = ENCODER (PC2)
livestream = NDIROUTER (Live Stream)
```

## Run
```
go build -o ndi-router
./ndi-router
```

## Update output via GET
```
curl http://localhost/updateOutput?output=main-screen&input=pc1
curl http://localhost/updateOutput?output=livestream&input=cam1
```

## Update output via WebSocket
```
let socket = new WebSocket("ws://localhost");

let emit = {
    Type: "UpdateOutput",
    Output: "main-screen",
    Input: "pc1",
}
socket.send(JSON.stringify(emit));

let emit = {
    Type: "UpdateOutput",
    Output: "livestream",
    Input: "cam1",
}
socket.send(JSON.stringify(emit));
```