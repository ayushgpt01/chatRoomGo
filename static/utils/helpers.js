import { LeftMessage } from "../components/LeftMessage.js";
import { RightMessage } from "../components/RightMessage.js";

export class Event {
  constructor(type, payload) {
    this.type = type;
    this.payload = payload;
  }
}

export function insertComponent(componentId, nodeId, attributes) {
  const component = document.createElement(componentId);
  for (const key in attributes) {
    component.setAttribute(key, attributes[key]);
  }
  document.getElementById(nodeId).appendChild(component.cloneNode(true));
}

export function addChatMessage({ senderType, message, id, username, status }) {
  insertComponent(
    senderType === "out"
      ? RightMessage.ComponentName
      : LeftMessage.ComponentName,
    "chat-window",
    {
      id,
      message,
      username,
      status,
    }
  );
}
