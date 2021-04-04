import Vue from 'vue'
import { SnotifyService } from 'vue-snotify/SnotifyService'
import { couchdb } from '~/plugins/couchdb'
import { DDB } from '~/plugins/DDB'

declare module 'vue/types/vue' {
  interface Vue {
    $snotify: SnotifyService
    $bus: Vue
    $couchdb: couchdb
    $ddb: DDB
  }
}
