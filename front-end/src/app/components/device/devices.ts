import { Device } from "./device";

/**
 * Devices has a Map all containing all devices with a key that is the same as the id.
 */
export class Devices {
  all: Map<string, Device>;

  constructor() {
    this.all = new Map<string, Device>();
  }

  /**
   * setDevice either updates an existing Device with the update methods or creates a new one.
   * @param jsonData json object with keys id, status and connection.
   */
  setDevice(jsonData) {
    if (this.all.has(jsonData.id)) {
      this.all.get(jsonData.id).updateStatus(jsonData.status);
      this.all.get(jsonData.id).updateConnection(jsonData.connection);
    } else {
      this.all.set(jsonData.id, new Device(jsonData));
    }
  }

  /**
   * Update the status of device front-end.
   * Update the component with id to status.
   */
  updateDevice(id, status) {
    if (this.all.has("front-end")) {
      const newStatus = {};
      newStatus[id] = status;
      this.all.get("front-end").updateStatus(newStatus);
    }
  }

  /**
   * getDevice is a getter for devices
   * @param dev device id
   */
  getDevice(dev: string) {
    if (this.all.has(dev)) {
      return this.all.get(dev);
    } else {
      return null;
    }
  }
}
