import { initWebSocket, sendEvent } from "./utils/websocket.js";

window.onload = function () {
  document
    .getElementById("sendMessageBtn")
    .addEventListener("click", handleSubmit);
  initWebSocket();
};

function handleSubmit(e) {
  e.preventDefault();
  const messageEle = document.getElementById("message");
  const message = messageEle.value;
  console.log("message", message);
  sendEvent("send_message", { message, senderType: "out" });
  messageEle.value = "";
}
