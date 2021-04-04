<template>
  <v-container>
    <v-progress-linear color="lime" indeterminate reverse :active="loading" />
    <dyn-table
      :dynObject="tableHeader"
      :dynObjectData="tableData"
      :addRowOnTop="tableNewRow"
    />
    <v-row no-gutters>
      <dyn-actions :actions="tableHeader.actions" @action="onAction" />
    </v-row>
  </v-container>
</template>

<script lang="ts">
import DynTable from "@/components/dynobject/Table.vue";
import { Vue, Component } from "vue-property-decorator";
import { DynObject } from "../../components/dynobject/DynObject";
import { DynObjectInput } from "~/components/dynobject/DynObjectInput";
import { DynObjectAction } from "~/components/dynobject/DynObjectAction";
import { qMessage } from "~/plugins/qommunicator";
import { DynObjectRow } from "~/components/dynobject/DynObjectRow";
import { $bus } from "~/plugins/eventBus";
import DynActions from "~/components/dynobject/Actions.vue";

export interface UserListElement {
  username: string;
}

@Component({
  components: { DynTable, DynActions },
})
export default class QmLog extends Vue {
  private loading?: boolean = false;
  private tableHeader?: DynObject = {} as DynObject;
  private tableData: DynObject = {} as DynObject;
  private tableNewRow?: DynObjectRow = {} as DynObjectRow;

  /*
  @Prop()
  public activated?: boolean;
*/

  mounted() {
    const newTableHeaderObject: DynObject = {} as DynObject;
    newTableHeaderObject.inputs = [
      {
        name: "direction",
        displayName: "Direction",
        type: "icon",
      } as DynObjectInput,
      {
        name: "topic",
        displayName: "Topic",
        type: "string",
      } as DynObjectInput,
      {
        name: "value",
        displayName: "Payload",
        type: "string",
      } as DynObjectInput,
    ];
    newTableHeaderObject.actions = [
      {
        name: "refresh",
        command: "log-refresh",
        response: "",
      } as DynObjectAction,
      {
        name: "ping",
        command: "log-ping",
        response: "",
      } as DynObjectAction,
    ];

    this.tableHeader = newTableHeaderObject;

    const $vm = this;
    $bus.$on("sse/qmessage", (message: qMessage) => {
      console.debug("New SSE-Message received, add it to log table");
      const newRow: DynObjectRow = {
        values: {
          direction: "mdi-help",
          topic: message.t,
          value: message.v,
        },
      } as DynObjectRow;

      if (message.directionSend === true) {
        newRow.values.direction = "mdi-export";
      } else {
        newRow.values.direction = "mdi-import";
      }

      this.tableNewRow = newRow;
    });

    $bus.$on("sse/qmessage/sended", (message: qMessage) => {
      console.debug("New SSE-Message sended, add it to log table");
      const newRow: DynObjectRow = {
        values: {
          direction: "mdi-help",
          topic: message.t,
          value: message.v,
        },
      } as DynObjectRow;

      if (message.directionSend === true) {
        newRow.values.direction = "mdi-export";
      } else {
        newRow.values.direction = "mdi-import";
      }

      this.tableNewRow = newRow;
    });

    this.onActionRefresh();
  }

  private onAction(action: string, payload: any) {
    switch (action) {
      case "log-ping":
        this.onActionPing();
        break;
      case "log-refresh":
        this.onActionRefresh();
        break;
    }
  }

  private onActionPing() {
    this.$qommunicator.sendToMqttRaw({
      t: "$SYS/ping",
    } as qMessage);
    return;
  }

  private onActionRefresh() {
    let tableData: DynObject = {} as DynObject;
    tableData.rows = [];

    for (let qm of this.$qommunicator.messageLog.toArray()) {
      const newRow: DynObjectRow = {
        values: {
          direction: "mdi-help",
          topic: qm.t,
          value: qm.v,
        },
      } as DynObjectRow;

      if (qm.directionSend === true) {
        newRow.values.direction = "mdi-export";
      } else {
        newRow.values.direction = "mdi-import";
      }

      // try to parse the json
      try {
        newRow.values.value = JSON.stringify(qm.v, undefined, "  ");
      } catch (e) {}

      tableData.rows.unshift(newRow);
    }
    this.tableData = tableData;
  }
}
</script>
