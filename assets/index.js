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
let counter = 0;

function fnv1aHash(data) {
  let hash = 0x811c9dc5;
  for (let i = 0; i < data.length; i++) {
    hash ^= data[i];
    hash = (hash >>> 0) * 0x01000193;
  }
  return hash >>> 0;
}

function encodeProtocol(command, counter, payload) {
  const payloadLength = payload.length;
  const buffer = new ArrayBuffer(1 + 2 + 4 + payloadLength + 4);
  const view = new DataView(buffer);
  let offset = 0;

  // A: command (opcode, 1 byte)
  view.setUint8(offset, command);
  offset += 1;

  // B: counter (2 bytes, big endian)
  view.setUint16(offset, counter, false);
  offset += 2;

  // C: payload length (4 bytes, big endian)
  view.setUint32(offset, payloadLength, false);
  offset += 4;

  // D: payload
  const payloadView = new Uint8Array(buffer, offset, payloadLength);
  payloadView.set(payload);
  offset += payloadLength;

  // F: checksum (FNV-1a, 32 bits, big endian)
  const checksum = fnv1aHash(new Uint8Array(buffer, 0, offset));
  view.setUint32(offset, checksum, false);
  offset += 4;

  return buffer;
}

function decodeProtocol(buffer) {
  if (buffer.length < 11) {
    throw new Error("Buffer muito curto " + buffer.length + ' "' + new TextDecoder().decode(buffer) + '"');
  }

  let offset = 0;

  // A: command (opcode, 1 byte)
  const command = buffer[offset];
  offset += 1;

  // B: counter (2 bytes, big endian)
  const counter = new DataView(buffer.slice(offset, offset + 2).buffer).getUint16(0, false);
  offset += 2;

  // C: payload length (4 bytes, big endian)
  const payloadLength = new DataView(buffer.slice(offset, offset + 4).buffer).getUint32(0, false);
  offset += 4;

  if (offset + payloadLength + 4 > buffer.length) {
    throw new Error("Comprimento de payload inválido");
  }

  // D: payload
  const payload = buffer.slice(offset, offset + payloadLength);
  offset += payloadLength;

  // F: checksum (FNV-1a, 32 bits, big endian)
  const checksum = new DataView(buffer.slice(offset, offset + 4).buffer).getUint32(0, false);
  // validate checksum

  return {
    command,
    counter,
    payloadLength,
    payload,
  };
}

function connectWS() {
  const { host, pathname: path, protocol: proto } = window.location;
  const url = `${proto === 'https:' ? 'wss' : 'ws'}://${host}${path === '/' ? '' : path}/ws`;
  ws = new WebSocket(url);
  ws.binaryType = 'arraybuffer';

  ws.onopen = () => {
    console.log('Conectado');
  };

  ws.onmessage = ({ data }) => {
    if (!data || data.byteLength === 0) {
      return;
    }

    let array = new Uint8Array(data);
    while (array.length > 0) {
      const { command, counter, payloadLength, payload } = decodeProtocol(array);
      switch (command) {
        case MSG:
          console.log(counter, new TextDecoder().decode(payload));
          break;
        case RUNJS: {
          const code = new TextDecoder().decode(payload);
          console.log(counter, code);
          try {
            const result = eval(code);
            console.log('Resultado:', result);
          } catch (e) {
            console.error('Erro:', e);
          }
          break;
        }
        case APPLYCSS: {
          const css = new TextDecoder().decode(payload);
          const style = document.createElement('style');
          style.textContent = css;
          document.head.appendChild(style);
          break;
        }
        case APPLYHTML: {
          const [id, ...htmlParts] = new TextDecoder().decode(payload).split('\n');
          const html = htmlParts.join('\n');
          const element = document.getElementById(id);
          if (element) {
            element.innerHTML = html;
          }
          break;
        }
        case LOADJS: {
          const urljs = new TextDecoder().decode(payload);
          const script = document.createElement('script');
          script.src = urljs;
          document.head.appendChild(script);
          break;
        }
        case LOADCSS: {
          const urlcss = new TextDecoder().decode(payload);
          const link = document.createElement('link');
          link.rel = 'stylesheet';
          link.type = 'text/css';
          link.href = urlcss;
          document.head.appendChild(link);
          break;
        }
        case LOADHTML: {
          const urlhtml = new TextDecoder().decode(payload);
          fetch(urlhtml)
            .then(response => response.text())
            .then(html => {
              const element = document.createElement('div');
              element.innerHTML = html;
              document.body.appendChild(element);
            });
          break;
        }
        case REGISTEREVENT: {
          const payloadText = new TextDecoder().decode(payload);
          const [eventType, label, id] = payloadText.split('\n');

          const element = document.getElementById(id);
          if (element) {
            element.addEventListener(eventType, (event) => {
              const eventData = {
                label,
                id,
                type: event.type,
                timestamp: event.timeStamp,
                // add more properties if needed
              };
              const payloadString = JSON.stringify(eventData);
              const payloadBytes = new TextEncoder().encode(payloadString);
              const messageBuffer = encodeProtocol(EVENT, 0, payloadBytes);
              ws.send(messageBuffer);
            });
          } else {
            console.warn(`element not found: ${id}`);
            // TODO: send error message to server
          }
          break;
        }
        default:
          console.log('Unknown command:',command, counter, payloadLength, payload);
          break;
      }
      if (array.length <= payloadLength + 11) {
        break;
      }
      array = array.slice(payloadLength + 11);
    }
  };

  ws.onerror = () => {
    console.log('error occurred, closing connection');
    ws.close();
  };

  ws.onclose = () => {
    console.log('reconnecting...');
    setTimeout(connectWS, 1000);
  };
}

// Connect to WebSocket server
window.onload = () => {
  console.log('connecting...');
  connectWS();
};

