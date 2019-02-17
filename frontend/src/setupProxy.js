const proxy = require("http-proxy-middleware")

module.exports = app => {
    app.use(proxy("/ws", {target: "http://10.0.3.6:8080", ws: true}))
}
