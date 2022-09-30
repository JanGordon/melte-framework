const exampleSocket = new WebSocket("ws://127.0.0.1:8080/hotReloadWS");

exampleSocket.onopen = (event) => {
    exampleSocket.send("Succesfully connected");
};

exampleSocket.onmessage = (event) => {
    window.location.reload()
    exampleSocket.send("reloaded")
}
console.log("connected")