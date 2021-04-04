<template>
  <v-container>
    <v-tabs v-model="tabIndex" grow>
      <v-tab v-for="item in tabs" :key="item.displayText" style="display: grid">
        <v-icon left>
          {{ item.icon }}
        </v-icon>
        {{ item.displayText }}
      </v-tab>

      <v-tabs-items v-model="tabIndex">
        <v-tab-item v-for="item in tabs" :key="item.displayText">
          <v-card flat>
            <component :is="item.component" v-if="item.active" />
          </v-card>
        </v-tab-item>
      </v-tabs-items>
    </v-tabs>
  </v-container>
</template>

<script lang="ts">
import { Vue, Component, Watch } from "vue-property-decorator";
import Users from "./users.vue";
import GopilotConnections from "./gopilot-connections.vue";
import SseConnections from "./sse-connections.vue";
import Log from "./log.vue";

export interface tabInfos {
  displayText: string;
  icon: string;
  active: boolean;
  component: string;
}

@Component({
  components: { Users, GopilotConnections, SseConnections, Log },
})
export default class LogoutPage extends Vue {
  private tabIndex: number = -1;
  private firstOneActive: boolean = false;

  private tabs: tabInfos[] = [
    {
      displayText: "Log",
      icon: "mdi-format-list-text",
      active: false,
      component: "log",
    },
    {
      displayText: "Users",
      icon: "mdi-account",
      active: false,
      component: "users",
    },
    /*
    {
      displayText: "SSE - Connections",
      icon: "mdi-lock",
      active: false,
      component: "sse-connections",
    },
    {
      displayText: "GP - Connections",
      icon: "mdi-lock",
      active: false,
      component: "gopilot-connections",
    },
    */
  ];

  @Watch("tabIndex")
  onTabChanged(value: number, oldValue: number) {
    if (oldValue > -1) {
      this.tabs[oldValue].active = false;
    }
    this.tabs[value].active = true;
  }
}
</script>
