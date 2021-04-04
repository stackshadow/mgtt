<template>
  <v-container>
    <v-progress-linear color="lime" indeterminate reverse :active="loading" />
    <dyn-table :dynObject="tableHeader" :dynObjectData="tableData" />
  </v-container>
</template>

<script lang="ts">
import DynTable from '@/components/dynobject/Table.vue'
import { Vue, Component } from 'vue-property-decorator'
import { DynObject } from '../../components/dynobject/DynObject'

@Component({
  components: { DynTable },
})
export default class SSEConnections extends Vue {
  private loading?: boolean = false
  private tableHeader?: DynObject = {} as DynObject
  private tableData?: DynObject = {} as DynObject

  /*
  @Prop()
  public activated?: boolean;
*/

  mounted() {
    // when a group-respond comes in, we save it
    this.$bus.$on('sse/client/table', this.onClientsTable)
    this.$bus.$on('sse/clients', this.onClients)
  }

  destroyed() {
    this.$bus.$off('sse/client/table', this.onClientsTable)
  }

  public onClientsTable(obj: DynObject) {
    this.tableHeader = undefined
    this.tableHeader = obj
  }

  public onClients(obj: DynObject) {
    this.tableData = undefined
    this.tableData = obj
  }

  /*
  @Watch("activated")
  onPropertyChanged(value: boolean, oldValue: boolean) {
    
  }
  */
}
</script>
