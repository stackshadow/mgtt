<template>
  <v-container> </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'nuxt-class-component'

@Component({
  components: {},
})
export default class LogoutPage extends Vue {
  public created() {
    this.$qommunicator.logout()
  }
}
</script>
