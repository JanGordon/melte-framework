try {
    var exampleSocket = new WebSocket("ws://"+ location.host +"/hotReloadWS");

    exampleSocket.onopen = (event) => {
        exampleSocket.send("Succesfully connected");
    };

    exampleSocket.onmessage = (event) => {
        setTimeout(() => { window.location.reload()}, 1000)
        console.log("REloading")
        exampleSocket.close()
        exampleSocket.send("reloaded")
    }
    console.log("connected")
} catch (e) {}