<template>
  <v-container>
    <form>
      <v-text-field
        v-model="username"
        :counter="30"
        label="Username"
        required
        outlined
      />
      <v-text-field
        v-model="password"
        label="Password"
        required
        outlined
        :type="showpassword ? 'text' : 'password'"
        :append-icon="showpassword ? 'mdi-eye' : 'mdi-eye-off'"
        @click:append="showpassword = !showpassword"
        @keyup.native.enter="submit"
      />

      <v-btn class="mr-4" @click="submit"> Login </v-btn>
    </form>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "nuxt-class-component";

@Component({
  components: {},
})
export default class ProjectPage extends Vue {
  private showpassword: boolean = false;
  private username: string = "";
  private password: string = "";

  private created() {}

  private async submit() {
    const result = await this.$qommunicator.authLogin(
      this.username,
      this.password
    );

    if (result === true) {
      this.$router.push("/mgtt");
    }
  }
}
</script>
