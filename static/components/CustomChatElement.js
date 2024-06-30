export class CustomChatElement extends HTMLElement {
  constructor(templateId) {
    super();
    this.templateId = templateId;
  }

  connectedCallback() {
    const template = document.getElementById(this.templateId).content;
    const clone = document.importNode(template, true);
    const spans = clone.querySelectorAll("span");
    const username = spans[0];
    const statusMessage = spans[1];
    const message = clone.querySelector("p");

    statusMessage.textContent = this.getAttribute("status");
    message.textContent = this.getAttribute("message");
    username.textContent = this.getAttribute("username");

    this.removeAttribute("status");
    this.removeAttribute("message");
    this.removeAttribute("username");

    this.appendChild(clone);
  }
}
