import { CustomChatElement } from "./CustomChatElement.js";

export class RightMessage extends CustomChatElement {
  static ComponentName = "right-msg-component";

  constructor() {
    super("right-message");
  }
}

customElements.define(RightMessage.ComponentName, RightMessage);
