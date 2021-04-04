<template>
  <vue-snotify />
</template>

<style>
.snotify {
  width: 600px;
}
</style>

<script lang="ts">
import Vue from 'vue'
import Component from 'nuxt-class-component'
import { Prop } from 'nuxt-property-decorator'
import Snotify, { SnotifyPosition, SnotifyToastConfig } from 'vue-snotify'
import { Notification } from '@/types/notification.ts'

// load the style
require('vue-snotify/styles/material.css')

const options = {
  toast: {
    position: SnotifyPosition.centerBottom,
  },
} as SnotifyToastConfig

// this will init snotify
const vuesnotify = Vue.use(Snotify, options)

@Component({
  components: { vuesnotify },
})
export default class snotify extends Vue {
  @Prop()
  public hideonprod?: boolean

  mounted() {
    const vm: Vue = this
    const $bus: Vue = vm.$bus

    $bus.$on('notify', (notification: Notification) => {
      if (notification.type === 'success') {
        vm.$snotify.success(notification.text, notification.title)
      }
      if (notification.type === 'warning') {
        vm.$snotify.warning(notification.text, notification.title)
      }
      if (notification.type === 'error') {
        vm.$snotify.error(notification.text, notification.title)
      }
    })
  }
}
</script>
