import Vue from "vue";

import axios, { AxiosResponse } from "axios";
import { RingBuffer } from "ring-buffer-ts";

import { $bus } from "./eventBus";
import { httpEndpoints } from "~/types/httpEndpoints";

export interface qMessage {
  o?: string; // SendAsMessage will send the message as json if its not ""
  r?: string; // ReqID
  s?: string; // timestamp
  t?: string; // topic
  tr?: string; // topic-reply
  v: any; // payload

  // from here we should not send over the channel
  send?: boolean; // this should not sended over channel
  status?: number;
  directionSend?: boolean; // true = outgoingSession
}

interface LoginCredentials {
  username: string;
  password?: string;
}

export class Qommunicator {
  private evtSource: EventSource | null = null;

  public messageLog: RingBuffer<qMessage>;

  constructor() {
    const $vm = this;
    this.messageLog = new RingBuffer<qMessage>(30);

    $bus.$on("mqtt/connected", (connected: boolean) => {
      if (connected === true) {
        this.connectToEvents();
      }
    });

    // emit the message over the bus
    // and log it to internal ring-buffer
    $bus.$on("sse/qmessage", (message: qMessage) => {
      // parse payload
      try {
        message.v = JSON.parse(message.v);
      } catch (e) {}

      // Add it to the log
      $vm.messageLog.add(message);

      // push it over the bus
      if (message.t !== undefined) {
        $bus.$emit(message.t, message.v);
      }
    });
  }

  public async post(path: string, payload: string): Promise<boolean> {
    let status: number = 0;

    await axios
      .post(path, payload)
      .catch(() => {
        return Promise.resolve(false);
      })
      .then((response: boolean | AxiosResponse) => {
        if (Object.prototype.hasOwnProperty.call(response, "status") === true) {
          if ((response as AxiosResponse).status === 200) {
            status = 200;
          } else {
            status = 400;
          }
        } else {
          status = 400;
        }
      });

    return Promise.resolve(status === 200);
  }

  public async sendToMqtt(message: qMessage) {
    message.o = "y";
    await this.post(
      httpEndpoints.SendQMessageToBackend,
      JSON.stringify(message)
    );
  }

  public async sendToMqttRaw(message: qMessage) {
    // Add it to the log
    let messageCopy = Object.assign({}, message);
    messageCopy.directionSend = true;
    this.messageLog.add(messageCopy);

    // delete elements that should not be sended
    message.send = undefined;
    message.status = undefined;
    message.directionSend = undefined;

    message.o = "";
    await this.post(
      httpEndpoints.SendQMessageToBackend,
      JSON.stringify(message)
    );

    $bus.$emit("sse/qmessage/sended", messageCopy);
  }

  // login will send an login request and store the authentication cookie
  public async authLogin(username: string, password: string): Promise<boolean> {
    const success = await this.post(
      httpEndpoints.Login,
      JSON.stringify({
        username,
        password,
      } as LoginCredentials)
    );

    if (success === false) {
      $bus.$emit("mqtt/connected", false);
      return Promise.resolve(false);
    }

    $bus.$emit("mqtt/connected", true);
    return Promise.resolve(true);
  }

  public connectToEvents(): Promise<boolean> {
    if (this.evtSource !== null) {
      // console.info('Already connected to SSE')
      return Promise.resolve(true);
    }

    this.evtSource = new EventSource(httpEndpoints.SSE);

    // we connected
    this.evtSource.onopen = function () {
      $bus.$emit("sse/connected", true);
    };

    // incoming message
    this.evtSource.onmessage = function (event) {
      const qMessage = JSON.parse(event.data) as qMessage | undefined;
      if (qMessage === undefined) {
        console.error("Can not parse qMessage from data");
        return;
      }

      $bus.$emit("sse/qmessage", qMessage);
    };

    // an error occurred
    this.evtSource.onerror = function (this: EventSource) {
      $bus.$emit("sse/connected", false);
    };

    return Promise.resolve(true);
  }

  public logout(): Promise<boolean> {
    if (this.evtSource !== null) {
      this.evtSource.close();
      this.evtSource = null;
    }

    $bus.$emit("sse/self/groups", ["all", "anonym"]);
    return Promise.resolve(true);
  }
}

declare module "vue/types/vue" {
  interface Vue {
    $qommunicator: Qommunicator;
  }
}

// and init the bus
export var $qommunicator = new Qommunicator();
Vue.prototype.$qommunicator = $qommunicator;
