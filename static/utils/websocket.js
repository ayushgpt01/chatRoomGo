import { Event, addChatMessage } from "./helpers.js";

export let conn = null;

export function initWebSocket() {
  if (window.WebSocket) {
    console.log("supports websockets");
    conn = new WebSocket("ws://" + document.location.host + "/ws");

    conn.onmessage = function (evt) {
      const eventData = JSON.parse(evt.data);
      const event = Object.assign(new Event(), eventData);
      console.log("event", event);
      routeEvent(event);
    };
  } else {
    alert("Not supporting websockets");
  }
}

function routeEvent(event) {
  if (!event.type) {
    alert("no 'type' field in event");
  }
  switch (event.type) {
    case "new_message":
      addChatMessage(event.payload);
      break;
    default:
      alert("unsupported message type");
      break;
  }
}

export function sendEvent(eventName, payload) {
  const event = new Event(eventName, payload);
  console.log(JSON.stringify(event));
  conn.send(JSON.stringify(event));
}
