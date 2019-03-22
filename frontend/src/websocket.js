class _WS {
    messageSubscriptions = {};
    modelAddedSubscriptions = {};

    ConnectToWebsocket() {
        if (!("WebSocket" in window)) {
            alert("Your browser doesn't support websocket, you need to update page yourself");
        } else {
            // TODO: connect to current domain
            const hostname = window.location.hostname;
            const port = window.location.protocol === "https:" ? 443 : window.location.port;
            const url = (port === 443 ? "wss://" : "ws://") + 
				hostname + (port === 433 || port === "80" ? "" : ":" + port) + "/ws";
            let ws = new WebSocket(url);
            ws.onopen = () => {
                console.log("websocket connected");
            };
            ws.onmessage = (event) => {
                this.handleMessage(event.data);
            };
            ws.onclose = () => {
                // on close
                console.log("websocket disconnected");
                setTimeout(() => this.ConnectToWebsocket(), 1000);
            };
            ws.onerror = (event) => {
				// alert("websocket error")
                console.log("websocket error");
                console.log(event);
				console.log(event.data);
				ws.close();
            }
        }
    }

    SubscribeOnMessage(messageType, callback) {
        if (typeof this.messageSubscriptions[messageType] === "undefined") {
            this.messageSubscriptions[messageType] = [];
        }
        this.messageSubscriptions[messageType].push(callback);
    }

    SubscribeOnModelAdded(modelName, callback) {
        if (typeof this.modelAddedSubscriptions[modelName] === "undefined") {
            this.modelAddedSubscriptions[modelName] = [];
        }
        this.modelAddedSubscriptions[modelName].push(callback);
    }

    handleMessage(message) {
        const jsonMessage = JSON.parse(message);
        const messageType = jsonMessage.type;

        console.log("got message of type '" + messageType + "' with content '" + jsonMessage + "'");

        this.notifyIfExist(this.messageSubscriptions, messageType, jsonMessage);

        switch (messageType) {
            case "model_added":
                const modelName = jsonMessage.data.model_name;
                const modelData = jsonMessage.data.data;
                this.notifyIfExist(this.modelAddedSubscriptions, modelName, modelData);
                break;
            default:
                console.log("unhandled message type " + messageType);
        }
    }

    notifyIfExist(map, key, value) {
        if (typeof map[key] !== "undefined") {
            for (let i in map[key]) {
                const callback = map[key][i];
                callback(value);
            }
        }
    }
}

let WS = new _WS();

export default WS;
