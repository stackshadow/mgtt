<template>
  <v-container>
    <simple-dialog type="delete" title="Delete user" />

    <!-- Editor -->
    <v-container v-if="editorVisible">
      <dyn-editor
        :mode="editorMode"
        :title="editorTitle"
        :primaryIcon="editorIcon"
        :inputs="tableHeader.inputs"
        :values.sync="editorValues"
        @action="onAction"
      />
    </v-container>

    <v-row no-gutters>
      <v-progress-linear color="lime" indeterminate reverse :active="loading" />
    </v-row>

    <!-- Table -->
    <v-row no-gutters>
      <v-col cols="12">
        <dyn-table
          :dynObject="tableHeader"
          :dynObjectData="tableData"
          @action="onAction"
        />
      </v-col>
    </v-row>

    <!-- Common actions -->
    <v-row no-gutters>
      <dyn-actions
        :actions="tableHeader.actions"
        :payload="editorValues"
        @action="onAction"
      />
    </v-row>
  </v-container>
</template>

<script lang="ts">
import DynTable from "@/components/dynobject/Table.vue";
import DynEditor from "@/components/dynobject/ValuesEditorDialog.vue";
import { Vue, Component } from "vue-property-decorator";
import { DynObject } from "../../components/dynobject/DynObject";
import { DynObjectInput } from "~/components/dynobject/DynObjectInput";
import { DynObjectAction } from "~/components/dynobject/DynObjectAction";
import { DynObjectRow } from "~/components/dynobject/DynObjectRow";
import { qMessage } from "~/plugins/qommunicator";
import DynActionButton from "~/components/dynobject/DynActionButton.vue";
import SimpleDialog from "@/components/core/simpledialog.vue";
import { UserListElement } from "~/plugins/mgtt";
import DynActions from "~/components/dynobject/Actions.vue";

const tableObject: DynObject = {
  inputs: [
    {
      name: "username",
      displayName: "Username",
      type: "string",
      valueMandatory: true,
      visible: true,
    } as DynObjectInput,
    {
      name: "password",
      displayName: "Password",
      type: "string",
      visible: false, // not visible in tables or overviews
      valueHidden: true,
    } as DynObjectInput,
    {
      name: "groups",
      displayName: "Goups",
      type: "[]string",
      hint: "Type a new group and press enter, or select one",
      visible: true, // not visible in tables or overviews
    } as DynObjectInput,
  ],
  actions: [
    {
      name: "new",
      command: "account-create-request",
      response: "",
    } as DynObjectAction,
    {
      name: "refresh",
      command: "account-refresh",
      response: "",
    } as DynObjectAction,
  ],
} as DynObject;

@Component({
  components: {
    DynTable,
    DynEditor,
    DynActions,
    DynActionButton,
    SimpleDialog,
  },
})
export default class QmUsers extends Vue {
  private loading?: boolean = false;
  private users?: UserListElement[];
  private tableHeader: DynObject = tableObject;
  private tableData?: DynObject = {} as DynObject;

  private editorVisible?: boolean = false;
  private editorMode?: string = "view";
  private editorTitle?: string = "View";
  private editorIcon?: string = "";
  private editorValues?: any = {};

  private onAction(action: string, payload: any) {
    switch (action) {
      case "account-create-request":
        this.onAccountCreateRequest();
        break;

      case "account-refresh":
        this.onAccountListRefresh();
        break;

      case "account-edit-request":
        this.onAccountEditRequest(payload.username as string);
        break;

      case "account-delete-request":
        this.onAccountDeleteRequest(payload.username as string);
        break;

      case "save":
        this.onAccountSave(this.editorValues as UserListElement);
        break;

      case "delete":
        this.onAccountDelete(this.editorValues as UserListElement);
        break;

      default:
        this.editorVisible = false;
    }
  }

  private onAccountCreateRequest() {
    this.editorMode = "create";
    this.editorTitle = "Create an new user";
    this.editorIcon = "mdi-account-plus";
    this.editorVisible = true;
  }

  private onAccountEditRequest(username: string) {
    this.$mgtt.Users.Get(username, (user: UserListElement) => {
      this.editorMode = "edit";
      this.editorTitle = "Edit " + user.username;
      this.editorValues = user;
      this.editorVisible = true;
    });
  }

  private onAccountSave(user: UserListElement) {
    this.$mgtt.Users.Save(user, (userSaved: UserListElement) => {
      // did we found it ?
      let foundIt: boolean = false;

      // find the user in the list and replace it
      this.users?.forEach(
        (myUser: UserListElement, index: number, array: UserListElement[]) => {
          if (myUser.username === userSaved.username) {
            array[index] = userSaved;
            foundIt = true;
          }
        }
      );

      if (foundIt === false) {
        this.users?.push(userSaved);
      }

      // and update the list
      this.userListFillTable();

      this.editorVisible = false;
    });
  }

  private onAccountDeleteRequest(username: string) {
    this.$mgtt.Users.Get(username, (user: UserListElement) => {
      this.editorMode = "delete";
      this.editorTitle = "Delete " + user.username;
      this.editorValues = user;
      this.editorVisible = true;
    });
  }

  private onAccountDelete(user: UserListElement) {
    this.$mgtt.Users.Delete(user.username, (saved: boolean) => {
      // find the user in the list and replace it
      this.users?.forEach(
        (myUser: UserListElement, index: number, array: UserListElement[]) => {
          if (myUser.username === user.username) {
            delete array[index];
          }
        }
      );

      // and update the list
      this.userListFillTable();

      this.editorVisible = false;
    });
  }

  private onAccountListRefresh() {
    // show the loading-bar
    this.loading = true;

    this.$mgtt.Users.List((userList: UserListElement[]) => {
      this.users = userList;
      this.userListFillTable();
    });
  }

  private userListFillTable() {
    const newTableDataObject: DynObject = {} as DynObject;
    newTableDataObject.rows = [];

    this.users?.forEach((user: UserListElement) => {
      const newRow: DynObjectRow = {
        values: {
          username: user.username,
          password: "",
          groups: user.groups,
        },
        actions: [
          {
            name: "edit",
            command: "account-edit-request",
            response: "",
          } as DynObjectAction,
          {
            name: "delete",
            command: "account-delete-request",
            response: "",
          } as DynObjectAction,
        ],
      } as DynObjectRow;

      newTableDataObject.rows?.push(newRow);
    });

    this.tableData = newTableDataObject;
    this.loading = false;
  }

  created() {
    this.onAccountListRefresh();
  }
}
</script>
