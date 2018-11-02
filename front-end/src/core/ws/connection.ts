
export async function createConnection(host: string, port: string): Promise<WebSocket> {
    var conn: WebSocket;
    var connectionString = "ws://" + host + ":" + port + "/ws";
    conn = new WebSocket(connectionString);
    console.info("Starting Websocket on : " + connectionString)
    conn.onclose = function (_) {
        console.info("Closed connection")
    };
    conn.onmessage = function (evt) {
        console.info("Message: " + evt.data)
    };
    return conn;
}