import Vue from "vue";

export interface Notification {
  type: string;
  title: string;
  text: string;
}

declare module "vue/types/vue" {
  interface Vue {
    $bus: Vue;
  }
}

// and init the bus
export const $bus = new Vue();
Vue.prototype.$bus = $bus;
