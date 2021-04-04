<template>
  <v-row justify="center">
    <v-dialog
      v-model="isShown"
      fullscreen
      hide-overlay
      transition="dialog-bottom-transition"
    >
      <v-card>
        <!-- The title -->
        <v-toolbar dark color="primary">
          <v-toolbar-title>{{ title }}</v-toolbar-title>
        </v-toolbar>

        <v-list-item>
          <!-- An image -->
          <v-list-item-avatar tile size="80"
            ><v-icon x-large v-if="PrimaryIcon !== ''">{{
              PrimaryIcon
            }}</v-icon>
          </v-list-item-avatar>

          <!-- Description -->
          <v-list-item-content>
            {{ description }}
          </v-list-item-content>
        </v-list-item>

        <component v-if="component" :is="component" />

        <!-- Actions -->
        <v-toolbar dark color="primary">
          <v-spacer></v-spacer>
          <v-btn
            v-for="button in buttons"
            :key="button.text"
            outlined
            dark
            @click="emitSignal(button.signal)"
          >
            <v-icon v-if="button.icon !== ''">{{ button.icon }}</v-icon>
            {{ button.text }}
          </v-btn>
        </v-toolbar>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'nuxt-class-component'
import { Prop, Watch } from 'nuxt-property-decorator'
import Snotify, { SnotifyPosition, SnotifyToastConfig } from 'vue-snotify'

interface Button {
  text: string
  icon: string
  signal: string
}

@Component({
  components: {}
})
export default class SimpleDialog extends Vue {
  @Prop() public visible?: boolean
  @Prop() public type?: string
  @Prop() public title?: string
  @Prop() public description?: string
  @Prop() private component?: string

  // Button 1
  @Prop() public yes?: boolean

  // Button 2
  @Prop() public no?: boolean
  @Prop() public save?: boolean

  // Button3
  @Prop() public cancel?: boolean

  public isShown?: boolean = false
  private PrimaryIcon?: string = ''
  private buttons?: Button[] = []

  created() {
    if (this.yes !== undefined || this.type === 'delete') {
      this.PrimaryIcon = 'mdi-alert-outline'
      this.buttons?.push({
        text: 'Yes',
        icon: 'mdi-hand-okay',
        signal: 'yes'
      } as Button)
    }

    if (this.no !== undefined || this.type === 'delete') {
      this.buttons?.push({
        text: '',
        icon: 'mdi-thumb-down-outline',
        signal: 'no'
      } as Button)
    }

    if (this.save !== undefined || this.type === 'save') {
      this.PrimaryIcon = 'mdi-help'
      this.buttons?.push({
        text: 'Save',
        icon: 'mdi-database-plus',
        signal: 'save'
      } as Button)
    }

    if (this.cancel !== undefined || this.type === 'save') {
      this.buttons?.push({
        text: 'Close',
        icon: 'mdi-close',
        signal: 'close'
      } as Button)
    }
  }

  @Watch('visible')
  onVisibleChanged(newValue?: boolean) {
    this.isShown = newValue
  }

  emitSignal(signal: string) {
    this.isShown = false
    this.$emit(signal)
  }
}
</script>
