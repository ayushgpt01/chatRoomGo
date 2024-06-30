import { CustomChatElement } from "./CustomChatElement.js";

export class LeftMessage extends CustomChatElement {
  static ComponentName = "left-msg-component";

  constructor() {
    super("left-message");
  }
}

customElements.define(LeftMessage.ComponentName, LeftMessage);
