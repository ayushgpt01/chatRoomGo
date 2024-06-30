const CHAT_WINDOW_ID = "chat-window";
/**
 * @type {WebSocket | null}
 */
let conn = null;

/**
 * Base class to create bubble element
 */
class CustomChatElement extends HTMLElement {
  templateId = "";

  constructor(id) {
    super();
    this.templateId = id;
  }

  connectedCallback() {
    // Create clone template
    const template = document.getElementById(this.templateId).content;
    const clone = document.importNode(template, true);

    // Get elements to update
    const spans = clone.querySelectorAll("span");
    const username = spans[0];
    const statusMessage = spans[1];
    const message = clone.querySelector("p");

    // Update element with attributes
    statusMessage.textContent = this.getAttribute("status");
    message.textContent = this.getAttribute("message");
    username.textContent = this.getAttribute("username");

    // Remove attributes
    this.removeAttribute("status");
    this.removeAttribute("message");
    this.removeAttribute("username");

    // Add to DOM
    this.appendChild(clone);
  }
}

/**
 * Class for creating Left Message Bubble Template for Chat window
 */
class LeftMessage extends CustomChatElement {
  static ComponentName = "left-msg-component";
  constructor() {
    super("left-message");
  }
}

/**
 * Class for creating Right Message Bubble Template for Chat window
 */
class RightMessage extends CustomChatElement {
  static ComponentName = "right-msg-component";
  constructor() {
    super("right-message");
  }
}

/**
 * Used to create new component and insert into DOM
 * @param {string} componentId
 * @param {string} nodeId
 * @param {Record<string,string>} attributes
 */
const insertComponent = (componentId, nodeId, attributes) => {
  const component = document.createElement(componentId);

  for (const key in attributes) {
    component.setAttribute(key, attributes[key]);
  }

  document.getElementById(nodeId).appendChild(component.cloneNode(true));
};

/**
 * Handler for adding message to chat window
 * @typedef MessageData
 * @property {string} id
 * @property {string} message
 * @property {"in" | "out"} senderType
 * @property {string} username
 *
 * @param {MessageData} messageData
 */
const addChatMessage = ({ senderType, message, id, username }) => {
  insertComponent(
    senderType === "out"
      ? RightMessage.ComponentName
      : LeftMessage.ComponentName,
    CHAT_WINDOW_ID,
    {
      id,
      message,
      username,
      status: "Delivered",
    }
  );
};

/**
 *
 * @param {MouseEvent} e
 */
const handleSubmit = (e) => {
  e.preventDefault();

  /**
   * @type {HTMLTextAreaElement}
   */
  const messageEle = document.getElementById("message");
  const message = messageEle.value;
  messageEle.value = "";
  sendEvent("send_message", message);
  // fetch("/addMessage", {
  //   method: "POST",
  //   body: JSON.stringify({ message }),
  // }).then((res) => {
  //   res.json().then((val) => {
  //     addChatMessage({ ...val, senderType: "out" });
  //   });
  // });
};

/**
 * Event is used to wrap all messages Send and Recieved
 * on the Websocket
 * The type is used as a RPC
 * */
class Event {
  // Each Event needs a Type
  // The payload is not required
  constructor(type, payload) {
    this.type = type;
    this.payload = payload;
  }
}

/**
 * routeEvent is a proxy function that routes
 * events into their correct Handler
 * based on the type field
 * */
function routeEvent(event) {
  if (event.type === undefined) {
    alert("no 'type' field in event");
  }
  switch (event.type) {
    case "new_message":
      console.log("new message");
      break;
    default:
      alert("unsupported message type");
      break;
  }
}

/**
 * sendEvent
 * eventname - the event name to send on
 * payload - the data payload
 * */
function sendEvent(eventName, payload) {
  // Create a event Object with a event named send_message
  const event = new Event(eventName, payload);
  // Format as JSON and send
  conn.send(JSON.stringify(event));
}

/**
 * Once the website loads, we want to apply listeners and connect to websocket
 * */
window.onload = function () {
  customElements.define(LeftMessage.ComponentName, LeftMessage);
  customElements.define(RightMessage.ComponentName, RightMessage);

  document
    .getElementById("sendMessageBtn")
    .addEventListener("click", handleSubmit);

  // Check if the browser supports WebSocket
  if (window["WebSocket"]) {
    console.log("supports websockets");
    // Connect to websocket
    conn = new WebSocket("ws://" + document.location.host + "/ws");

    // Add a listener to the onmessage event
    conn.onmessage = function (evt) {
      console.log(evt);
      // parse websocket message as JSON
      const eventData = JSON.parse(evt.data);
      // Assign JSON data to new Event Object
      const event = Object.assign(new Event(), eventData);
      // Let router manage message
      routeEvent(event);
    };
  } else {
    alert("Not supporting websockets");
  }
};
