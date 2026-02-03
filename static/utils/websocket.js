import { Event, addChatMessage } from "./helpers.js";

export let conn = null;

export function initWebSocket() {
  if (window.WebSocket) {
    console.log("supports websockets");
    conn = new WebSocket("ws://" + document.location.host + "/ws");

    conn.onopen = function () {
      console.log("WebSocket connection established");
    };

    conn.onmessage = function (evt) {
      const eventData = JSON.parse(evt.data);
      const event = Object.assign(new Event(), eventData);
      console.log("event", event);
      routeEvent(event);
    };

    conn.onclose = function () {
      console.log("WebSocket connection closed, attempting to reconnect...");
      setTimeout(initWebSocket, 5000); // Reconnect after 5 seconds
    };

    conn.onerror = function (error) {
      console.error("WebSocket error:", error);
    };
  } else {
    alert("WebSocket not supported by your browser.");
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
  if (conn && conn.readyState === WebSocket.OPEN) {
    const event = new Event(eventName, payload);
    conn.send(JSON.stringify(event));
  } else {
    console.log("WebSocket connection is not open");
  }
}
