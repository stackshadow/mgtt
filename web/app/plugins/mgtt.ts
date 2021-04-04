import Vue from "vue";
import { $bus } from "./eventBus";
import { $qommunicator, qMessage, Qommunicator } from "./qommunicator";

export class MGTT {
  public Users: UsersClass = new UsersClass();
  public Groups: GroupsClass = new GroupsClass();
}

export interface UserListElement {
  username: string;
  password?: string;
  groups?: string[];
}

class UsersClass {
  public GetMyUserName(callback?: (username: string) => void) {
    $qommunicator.sendToMqttRaw({
      t: "$SYS/self/user/get",
      v: "",
    } as qMessage);

    if (callback !== undefined) {
      $bus.$once("$SYS/self/user/json", (userJson: UserListElement) => {
        callback(userJson.username);
      });
    }
  }

  public List(callback: (userList: UserListElement[]) => void) {
    $qommunicator.sendToMqttRaw({
      t: "$SYS/auth/users/list/get",
      v: "",
    } as qMessage);

    $bus.$once("$SYS/auth/users/list/json", (users: UserListElement[]) => {
      callback(users);
    });
  }

  public Get(username: string, callback: (user: UserListElement) => void) {
    $qommunicator.sendToMqttRaw({
      t: "$SYS/auth/user/" + username + "/get",
      v: "",
    } as qMessage);

    $bus.$once(
      "$SYS/auth/user/" + username + "/json",
      (user: UserListElement) => {
        callback(user);
      }
    );
  }

  public Save(
    changedUser: UserListElement,
    callback: (user: UserListElement) => void
  ) {
    $qommunicator.sendToMqttRaw({
      t: "$SYS/auth/user/" + changedUser.username + "/set",
      v: JSON.stringify(changedUser),
    } as qMessage);

    $bus.$once(
      "$SYS/auth/user/" + changedUser.username + "/set/success",
      (user: UserListElement) => {
        callback(user);
      }
    );
  }

  public Delete(username: string, callback: (deleted: boolean) => void) {
    $qommunicator.sendToMqttRaw({
      t: "$SYS/auth/user/" + username + "/delete",
      v: "",
    } as qMessage);

    $bus.$once(
      "$SYS/auth/user/" + username + "/delete/success",
      (success: boolean) => {
        callback(success);
      }
    );
  }
}

class GroupsClass {
  public GetMyGroups(callback?: (groups: string[]) => void) {
    $qommunicator.sendToMqttRaw({
      t: "$SYS/self/groups/get",
      v: "",
    } as qMessage);

    if (callback !== undefined) {
      $bus.$once("$SYS/self/groups/json", (groups: string[]) => {
        callback(groups);
      });
    }
  }
}

declare module "vue/types/vue" {
  interface Vue {
    $mgtt: MGTT;
  }
}

// and init the bus
Vue.prototype.$mgtt = new MGTT();
