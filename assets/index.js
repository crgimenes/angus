const MSG = 0x1; // message to client console (debug)
const RUNJS = 0x2; // run JavaScript code
const APPLYCSS = 0x3; // apply CSS to the page
const APPLYHTML = 0x4; // apply HTML to the page
const LOADJS = 0x5; // load JavaScript file from URL
const LOADCSS = 0x6; // load CSS file from URL
const LOADHTML = 0x7; // load HTML file from URL
const REGISTEREVENT = 0x8; // register event listener on element
const EVENT = 0x9; // event from client

let ws;

function encodeProtocol(command, payload) {
  const payloadLength = payload.length;
  const buffer = new ArrayBuffer(5 + payloadLength);
  const view = new DataView(buffer);

  view.setUint8(0, command);
  view.setUint32(1, payloadLength, false);
  const payloadView = new Uint8Array(buffer, 5, payloadLength);
  payloadView.set(payload);

  return buffer;
}

function decodeProtocol(buffer) {
  if (buffer.length < 5) {
    throw new Error("Buffer too short " + buffer.length);
  }

  const command = buffer[0];
  const payloadLength = new DataView(buffer.slice(1, 5).buffer).getUint32(
    0,
    false,
  );
  const payload = buffer.slice(5, 5 + payloadLength);

  return {
    command,
    payloadLength,
    payload,
  };
}

function connectWS() {
  const { host, pathname: path, protocol: proto } = window.location;
  const url = `${proto === "https:" ? "wss" : "ws"}://${host}${path === "/" ? "" : path}/ws`;
  ws = new WebSocket(url);
  ws.binaryType = "arraybuffer";

  ws.onopen = () => {
    console.log("Connected to server");
  };

  ws.onmessage = ({ data }) => {
    if (!data || data.byteLength === 0) {
      return;
    }

    let array = new Uint8Array(data);
    while (array.length > 0) {
      const { command, payloadLength, payload } = decodeProtocol(array);
      switch (command) {
        case MSG:
          console.log(new TextDecoder().decode(payload));
          break;
        case RUNJS: {
          const code = new TextDecoder().decode(payload);
          try {
            const result = eval(code);
            console.log("Resultado:", result);
          } catch (e) {
            console.error("Erro:", e);
          }
          break;
        }
        case APPLYCSS: {
          const css = new TextDecoder().decode(payload);
          const style = document.createElement("style");
          style.textContent = css;
          document.head.appendChild(style);
          break;
        }
        case APPLYHTML: {
          const [id, ...htmlParts] = new TextDecoder()
            .decode(payload)
            .split("\n");
          const html = htmlParts.join("\n");
          const element = document.getElementById(id);
          if (element) {
            element.innerHTML = html;
          }
          break;
        }
        case LOADJS: {
          const urljs = new TextDecoder().decode(payload);
          const script = document.createElement("script");
          script.src = urljs;
          document.head.appendChild(script);
          break;
        }
        case LOADCSS: {
          const urlcss = new TextDecoder().decode(payload);
          const link = document.createElement("link");
          link.rel = "stylesheet";
          link.type = "text/css";
          link.href = urlcss;
          document.head.appendChild(link);
          break;
        }
        case LOADHTML: {
          const urlhtml = new TextDecoder().decode(payload);
          fetch(urlhtml)
            .then((response) => response.text())
            .then((html) => {
              const element = document.createElement("div");
              element.innerHTML = html;
              document.body.appendChild(element);
            });
          break;
        }
        case REGISTEREVENT: {
          const payloadText = new TextDecoder().decode(payload);
          const [eventType, label, id] = payloadText.split("\n");

          const element = document.getElementById(id);
          if (element) {
            element.addEventListener(eventType, (event) => {
              const eventData = {
                label,
                id,
                type: event.type,
                // add more properties if needed
              };
              const payloadString = JSON.stringify(eventData);
              const payloadBytes = new TextEncoder().encode(payloadString);
              const messageBuffer = encodeProtocol(EVENT, payloadBytes);
              ws.send(messageBuffer);
            });
          } else {
            console.warn(`element not found: ${id}`);
            // TODO: send error message to server
          }
          break;
        }
        default:
          console.log("Unknown command:", command, payloadLength, payload);
          break;
      }
      if (array.length <= payloadLength + 5) {
        break;
      }
      array = array.slice(payloadLength + 5);
    }
  };

  ws.onerror = () => {
    console.log("error occurred, closing connection");
    ws.close();
  };

  ws.onclose = () => {
    console.log("reconnecting...");
    setTimeout(connectWS, 1000);
  };
}

// Connect to WebSocket server
window.onload = () => {
  console.log("connecting...");
  connectWS();
};

// set bootstrap dark mode or light mode based browser preference
document.documentElement.setAttribute(
  "data-bs-theme",
  window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
);
