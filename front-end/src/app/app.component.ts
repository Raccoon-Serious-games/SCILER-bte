import { Component, OnDestroy, OnInit, ViewEncapsulation } from "@angular/core";
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Message } from "./message";
import { JsonConvert } from "json2typescript";
import { MatSnackBar, MatSnackBarConfig } from "@angular/material";
import { Observable, Subscription, timer } from "rxjs";
import { Devices } from "./components/device/devices";
import { Puzzles } from "./components/puzzle/puzzles";
import { Timers } from "./components/timer/timers";
import { Logger } from "./logger";
import { Camera } from "./camera/camera";
import { Hint } from "./components/hint/hint";
import { formatMS, formatTime } from "./components/timer/timer";
import { FullScreen } from "./fullscreen";

/**
 * This is the main application, controlling all actions that can happen.
 * It keeps track of the main data objects and communicates to the back-end.
 */
@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css", "../assets/css/main.css"],
  encapsulation: ViewEncapsulation.None
})
export class AppComponent extends FullScreen implements OnInit, OnDestroy {
  // Variables for the home screen
  title = "SCILER";
  nameOfRoom = "Super awesome escape";

  // Necessary tools
  jsonConvert: JsonConvert;
  logger: Logger;
  subscription: Subscription;

  // Keeping track of data
  deviceList: Devices;
  puzzleList: Puzzles;
  hintList: Hint[];
  configErrorList: string[];
  cameras: Camera[];
  selectedCamera: string;
  selectedCamera2: string;
  openSecondCamera = false;
  timerList: Timers;
  displayTime: string;
  everySecond: Observable<number> = timer(0, 1000);

  /**
   * When starting the application the first time, inject the parameters.
   * Initialize all the attributes of the application, subscribe to the topics of the broker,
   * and ask for the set-up of the back-end.
   *
   * @param mqttService for communication with back-end
   * @param snackBar material design message pop-up framework
   */
  constructor(private mqttService: MqttService, private snackBar: MatSnackBar) {
    super();
    this.logger = new Logger();
    this.jsonConvert = new JsonConvert();
    this.initializeVariables();

    const topics = ["front-end"];
    for (const topic of topics) {
      this.subscribeNewTopic(topic);
    }

    this.mqttService.onConnect.subscribe(() => {
      this.logger.log("info", "Connected to broker");
      this.sendInstruction([{ instruction: "send setup" }]);
      this.sendConnection(true);
      this.initializeTimers();
    });

    this.mqttService.onOffline.subscribe(() => {
      this.logger.log("error", "Connection to broker lost");
      this.setConnectionAllDevices(false);
    });
  }

  /**
   * Sets connection of all devices, starting as false, until message received telling it's connected.
   */
  private setConnectionAllDevices(connection: boolean) {
    for (const tuple of this.deviceList.all) {
      const device = tuple[1];
      device.connection = false;
    }
  }

  ngOnInit(): void {}

  /**
   * Set all the variables to their default state, removing old data.
   * Set the duration timer for the escape room to 0, this will be updated when data is received from back-end.
   */
  initializeVariables() {
    this.deviceList = new Devices();
    this.puzzleList = new Puzzles();
    this.hintList = [];
    this.configErrorList = [];
    this.cameras = [];
    this.timerList = new Timers();
    const generalTimer = { id: "general", duration: 0, state: "stateIdle" };
    this.timerList.setTimer(generalTimer);
  }

  /**
   * When the user leaves the app, tell the back-end about the disconnect.
   * Then, the broker subscriptions should be cleaned up
   * and the connection with the broker closed.
   */
  ngOnDestroy(): void {
    this.sendConnection(false);
    this.mqttService.disconnect();
  }

  /**
   * Subscribe to a certain topic from the broker.
   * Also, tell the subscription to process when a message is received on that topic.
   */
  private subscribeNewTopic(topic: string): void {
    this.subscription = this.mqttService
      .observe(topic)
      .subscribe((message: IMqttMessage) => {
        this.logger.log(
          "info",
          "received on topic " +
            message.topic +
            ", message: " +
            message.payload.toString()
        );
        this.processMessage(message.payload.toString());
      });
    this.logger.log("info", "subscribed to topic: " + topic);
  }

  /**
   * Send an instruction to the broker, over topic `back-end`.
   * @param instructions to be sent.
   */
  public sendInstruction(instructions: any[]) {
    const msg = new Message(
      "front-end",
      "instruction",
      new Date(),
      instructions
    );
    let jsonMessage: string = JSON.stringify(this.jsonConvert.serialize(msg));
    this.mqttService.unsafePublish("back-end", jsonMessage);
    for (const inst of instructions) {
      if ("config" in inst) {
        msg.contents = { config: "contents to long to print" };
        jsonMessage = JSON.stringify(this.jsonConvert.serialize(msg));
      }
    }
    this.logger.log("info", "sent instruction message: " + jsonMessage);
  }

  /**
   * Send a status to the broker, over topic `back-end`.
   * @param start start status to be sent.
   * @param stop stop status to be sent.
   */
  public sendStatus(start, stop) {
    const msg = new Message("front-end", "status", new Date(), {
      start,
      stop
    });
    const jsonMessage: string = this.jsonConvert.serialize(msg);
    this.mqttService.unsafePublish("back-end", JSON.stringify(jsonMessage));
    this.logger.log(
      "info",
      "sent status message: " + JSON.stringify(jsonMessage)
    );
  }

  /**
   * Send a connection update to the broker, over topic `back-end`.
   * @param connected connection status to be sent.
   */
  public sendConnection(connected: boolean) {
    const msg = new Message("front-end", "connection", new Date(), {
      connection: connected
    });
    const jsonMessage: string = this.jsonConvert.serialize(msg);
    this.mqttService.unsafePublish("back-end", JSON.stringify(jsonMessage));
    this.logger.log(
      "info",
      "sent connection message: " + JSON.stringify(jsonMessage)
    );
  }

  /**
   * Process the incoming message, depending on its type.
   * @param jsonMessage json message received.
   */
  private processMessage(jsonMessage: string) {
    const msg: Message = Message.deserialize(jsonMessage);

    switch (msg.type) {
      case "confirmation": {
        this.processConfirmation(msg);
        break;
      }
      case "instruction": {
        this.processInstruction(msg.contents);
        break;
      }
      case "status": {
        this.deviceList.setDevice(msg.contents);

        // When the back-end/front-end disconnects, all devices are disconnected
        if (msg.contents.id === "front-end" && !msg.contents.connection) {
          this.setConnectionAllDevices(false);
        }
        break;
      }
      case "event status": {
        this.puzzleList.updatePuzzles(msg.contents);
        break;
      }
      case "time": {
        this.timerList.setTimer(msg.contents);
        break;
      }
      case "setup": {
        this.processSetUp(msg.contents);
        break;
      }
      // when a config is checked by the back-end it returns a list of found errors, these should be displayed
      case "config": {
        this.configErrorList = msg.contents.errors;
        break;
      }
      // when a config has be checked and put to use (only possible on no errors), notify the user
      case "new config": {
        this.openSnackbar("using new config: " + msg.contents.name, "");
        break;
      }
      default:
        this.logger.log("error", "received invalid message type " + msg.type);
        break;
    }
  }

  /**
   * When the front-end receives confirmation message from client computer
   * that an instruction was completed, show the message to the user.
   */
  private processConfirmation(jsonData) {
    for (const instruction of jsonData.contents.instructed.contents) {
      const display =
        "received confirmation from " +
        jsonData.deviceId +
        " for instruction: " +
        instruction.instruction;
      this.openSnackbar(display, "");
    }
  }

  /**
   * Process instruction messages. The types that exist:
   * reset - reset the front-end's device status
   * status update - send front-end's connection status to back-end
   * test - perform a test on the front-end
   */
  private processInstruction(jsonData) {
    for (const action of jsonData) {
      switch (action.instruction) {
        case "reset":
          {
            this.deviceList.setDevice({
              id: "front-end",
              connection: true,
              status: {
                start: 0,
                stop: 0
              }
            });
          }
          break;
        case "status update": {
          this.sendConnection(true);
          break;
        }
        case "test": {
          this.openSnackbar("performing instruction test", "");
          break;
        }
      }
    }
  }

  /**
   * The setup contain:
   * the name of the room to display in app
   * the camera links to select in camera view
   * the rule descriptions for in the puzzle table
   * the map with hints per puzzle to display in hint selection box
   *
   * @param jsonData with name, camera, events, hints
   */
  private processSetUp(jsonData) {
    this.nameOfRoom = jsonData.name;

    const cameraData = jsonData.cameras;
    this.cameras = [];
    if (cameraData !== null) {
      for (const cam of cameraData) {
        this.cameras.push(new Camera(cam));
      }
    }

    const rules = jsonData.events;
    this.puzzleList = new Puzzles();
    for (const rule in rules) {
      if (rules.hasOwnProperty(rule)) {
        this.puzzleList.addPuzzle(rule, rules[rule]);
      }
    }

    const allHints = jsonData.hints;
    this.hintList = [];
    for (const puzzle in allHints) {
      if (allHints.hasOwnProperty(puzzle)) {
        const hints = [];
        for (const index in allHints[puzzle]) {
          if (allHints[puzzle].hasOwnProperty(index)) {
            hints.push(allHints[puzzle][index]);
          }
        }
        this.hintList.push(new Hint(puzzle, hints));
      }
    }
  }

  /**
   * Initialize the timers to listen to everySecond and set their state accordingly.
   */
  private initializeTimers() {
    this.subscription = this.everySecond.subscribe(seconds => {
      for (const aTimer of this.timerList.getAll().values()) {
        if (aTimer.state === "stateActive") {
          aTimer.tick();
        }
        if (aTimer.duration <= 0) {
          aTimer.state = "stateIdle";
        }
      }
      this.displayTime = formatMS(
        this.timerList.getTimer("general").getTimeLeft()
      );
    });
  }

  /**
   * Before using new configuration, first stop the current timer subscription.
   * Otherwise time runs double.
   */
  private stopTimers() {
    this.subscription.unsubscribe();
  }

  /**
   * Opens snackbar with duration of 3 seconds.
   * @param message displays this message
   * @param action: button to display - optional use
   */
  public openSnackbar(message: string, action: string) {
    const config = new MatSnackBarConfig();
    config.duration = 3000;
    config.panelClass = ["custom-snack-bar"];
    this.snackBar.open(message, action, config);
  }

  /**
   * Return the current time to display.
   */
  getCurrentTime() {
    const date = new Date();
    return formatTime(date.getTime(), date.getTimezoneOffset());
  }

  /**
   * Stops timers, then create new variables and timers
   */
  public resetConfig() {
    this.stopTimers();
    this.initializeVariables();
    this.initializeTimers();
  }
}
