<template>
  <v-container>
    <v-progress-linear color="lime" indeterminate reverse :active="loading" />
  </v-container>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";

@Component({
  components: {},
})
export default class GopilotConnections extends Vue {
  private loading: boolean = false;

  /*
  @Prop()
  public activated?: boolean;
*/

  mounted() {
    // console.log(this.loading);
  }

  /*
  @Watch("activated")
  onPropertyChanged(value: boolean, oldValue: boolean) {
    
  }
  */
}
</script>
