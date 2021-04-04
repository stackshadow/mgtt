<template>
  <v-app>
    <div class="v-application d-flex flex-row" data-app="true">
      <snotify style="z-index: 2000" />

      <v-container fluid style="margin: 0px; padding: 0px; width: 100%">
        <v-row no-gutters>
          <v-app-bar
            color="accent-4"
            dense
            class="flex-col flex-grow-1"
            cols="12"
          >
            <v-toolbar-title>MGTT</v-toolbar-title>

            <v-spacer />
            <v-btn
              icon
              :color="$vuetify.theme.dark ? 'yellow' : 'dark'"
              @click="$vuetify.theme.dark = !$vuetify.theme.dark"
              ><v-icon v-if="$vuetify.theme.dark === true">
                mdi-weather-sunny
              </v-icon>
              <v-icon v-if="$vuetify.theme.dark === false">
                mdi-moon-waning-crescent
              </v-icon>
            </v-btn>

            <v-chip class="ma-2" color="primary" label>
              <v-icon left> mdi-account-circle-outline </v-icon>
              {{ userName }}
            </v-chip>
            <v-btn
              v-if="isLoggedIn === false || userName === '_anonym'"
              icon
              link
              nuxt
              to="/user/login"
            >
              <v-icon> mdi-login </v-icon>
            </v-btn>
            <v-btn
              v-if="isLoggedIn === true && userName !== '_anonym'"
              icon
              link
              nuxt
              color="orange"
              to="/user/logout"
            >
              <v-icon> mdi-logout </v-icon>
            </v-btn>
          </v-app-bar>
        </v-row>
        <v-row no-gutters>
          <v-col no-gutters style="margin: 0px; padding: 0px; width: 100%">
            <nuxt />
          </v-col>
        </v-row>
      </v-container>
    </div>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "nuxt-class-component";

import snotify from "@/components/core/snotify.vue";

export interface qommunicatorMessage {
  s?: string; // timestamp
  t: string; // topic
  tr?: string; // topic-reply
  v: any; // payload
}

@Component({
  components: {
    snotify,
  },
})
export default class ProjectPage extends Vue {
  private showmenu: boolean = false;
  private currentUserRoles: string[] = ["all", "anonym", "debug"];

  private isLoggedIn: boolean = false;

  private userName: string = "";
  private userGroups: string[] = [];

  public page?: Object = {
    title: "",
  };

  async mounted() {
    const $vm: ProjectPage = this;

    this.$bus.$on("mqtt/connected", this.onMQTTConnected);

    // we login to mqtt without an user ( anonymouse )
    let result = await this.$qommunicator.authLogin("", "");
    if (result === true) {
      this.$router.push("/mgtt");
    }
  }

  // Declared as computed property getter
  togglemenu() {
    if (this.showmenu === false) {
      this.showmenu = true;
    } else {
      this.showmenu = false;
    }
  }

  private onMQTTConnected(connected: boolean) {
    if (connected === true) {
      this.$snotify.success("Erfolgreich", "Login");

      // get my-username my-group
      this.$mgtt.Users.GetMyUserName((username: string) => {
        this.userName = username;
        if (this.userName !== "") {
          this.isLoggedIn = true;
        }
      });
      this.$mgtt.Groups.GetMyGroups();
    } else {
      this.$snotify.error("Fehlgeschlagen", "Login");
      this.isLoggedIn = false;
    }
  }
}
</script>

<style scoped>
.theme--dark {
  background-color: var(--v-background-base, #121212) !important;
}

.theme--dark.v-data-table {
  background-color: #494949 !important;
}

.theme--dark.v-application {
  background-color: var(--v-background-base, #121212) !important;
}
.theme--light.v-application {
  background-color: var(--v-background-base, white) !important;
}
</style>
