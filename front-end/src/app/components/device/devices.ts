import { Device } from "./device";

/**
 * Class keeping track of the devices, through use of map with device id's.
 */
export class Devices {
  all: Map<string, Device>;

  constructor() {
    this.all = new Map<string, Device>();
  }

  /**
   * Set the new status and connection state of a certain device.
   * If the device did not yet exist, create a new one.
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
   * Return device with id dev.
   */
  getDevice(dev: string) {
    if (this.all.has(dev)) {
      return this.all.get(dev);
    } else {
      return null;
    }
  }
}
