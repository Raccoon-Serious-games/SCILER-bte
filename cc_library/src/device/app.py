from datetime import datetime
import json

import paho.mqtt.client as mqtt


def on_connect(client, userdata, flags, rc):
    """
    When trying to connect to the broker,
    on_connect will return the result of this action.
    """
    if rc == 0:
        client.connected_flag = True  # set flag
        print("connected OK")
    else:
        print("Bad connection Returned code=", rc)
        client.bad_connection_flag = True


def on_disconnect(client, userdata, rc):
    """
    When disconnecting from the broker, on_disconnect prints the reason.
    """
    msg_dict = {
        "device_id": "library",
        "time_sent": datetime.now().strftime("%d-%m-%YT%H:%M:%S"),
        "type": "connection",
        "message": {"connection": False},
    }

    msg = json.dumps(msg_dict)
    client.publish("connection", msg)
    print("disconnecting reason  " + str(rc))
    client.connected_flag = False
    client.disconnect_flag = True


def on_log(client, userdata, level, buf):
    """
    Very annoying logger that logs everything happening with the mqtt client.
    """
    print("log: ", buf)


class App:
    """
    Class App sets up the connection and the right handler
    """

    def __init__(self, config, device):
        self.device = device
        self.config = json.load(config)
        self.name = self.config.get("id")
        self.info = self.config.get("info")
        self.host = self.config.get("host")
        self.client = mqtt.Client(self.name)
        self.client.on_message = self.on_message
        self.client.on_log = on_log
        self.client.on_connect = on_connect
        self.client.on_disconnect = on_disconnect

    def start(self):
        """
        Starting method to call from the starting script.
        """
        self.connect()

    def subscribe_topic(self, topic):
        """
        Method to call to subscribe to a topic the
        device wants to recieve from the broker.
        """
        print("subscribed to topic", topic)
        self.client.subscribe(topic=topic)

    def send_status_message(self, msg):
        """
        Method to send status messages to the topic status.
        msg should be a dictionary/json with components
         as keys and its status as value
        """
        jsonmsg = {
            "device_id": self.name,
            "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
            "type": "status",
            "contents": eval(msg),
        }
        msg = json.dumps(jsonmsg)
        print(str(msg))
        self.client.publish("status", msg)

    def on_message(self, client, userdata, message):
        """
        This method is called when the client receives
         a message from the broken for a subscribed topic.
        The message is printed and send through to the handler.
        """
        print("message received ", str(message.payload.decode("utf-8")))
        print("message topic=", message.topic)
        self.handle(message)

    def connect(self):
        """
        Connect method to set up the connection to the broker.
        When connected:
        sends message to topic connection to say its connected,
        subscribes to topic "test"
        starts loop_forever
        """
        while True:
            try:
                self.client.connect(self.host, 1883, keepalive=60)
                print("we zijn connected")
                msg_dict = {
                    "device_id": self.name,
                    "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
                    "type": "connection",
                    "message": {"connection": True},
                }

                msg = json.dumps(msg_dict)
                self.client.publish("connection", msg)
                self.subscribe_topic("test")
                self.client.loop_forever()
                break
            except:
                print("alles is kapot")

    def handle(self, message):
        """
        Interpreter of incoming messages.
        Correct device mapper is called with the content of the message.
        """
        message = message.payload.decode("utf-8")
        message = json.loads(message)
        self.device.incoming_instruction(message.get("contents"))
