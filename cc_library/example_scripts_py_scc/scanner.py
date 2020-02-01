import os
from sciler.device import Device
from evdev import InputDevice, categorize, ecodes

try:
    from evdev import InputDevice, categorize, ecodes
except ImportError:
    print("Error import EVDEV. Module Room_door_handler can't be used")

import asyncio


class Scanner(Device):

    def get_status(self):
        """
        Returns status of all custom components, in json format.90311031

        """
        return {
            "code": self.code,
            "scanned": [self.scanned[len(self.scanned) - 2], self.scanned[len(self.scanned) - 1]]
        }

    def perform_instruction(self, contents):
        """
        Defines how instructions are handled,
        for all instructions defined in output of device in config.
        :param contents: contains instruction tag and calls the appropriate functions.
        :return boolean: True if instruction was valid and False if illegal instruction
        was sent or error occurred such that instruction could not be performed.
        Returns tuple, with boolean and None if True and the failed action if false.
        """

    def test(self):
        """
        Defines test sequence for device.
        """
        self.log("test")

    def reset(self):
        """
        Defines a reset sequence for device.
        """
        self.code = 0
        self.scanned = [0, 0]
        self.log("reset")

    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "scanner.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.code = 0
        self.scanned = [0, 0]

    def main(self):
        """
        This method should be overridden in the subclass.
        It should initialize the SccLib class with config file and device class.
        It should also add event listeners to GPIO for all input components.
        """

        def async_loop():
            while True:
                # self.code = eval(input())
                self.code = eval(read_barcode("/dev/input/event0"))
                print(self.code)
                self.scanned.append(self.code)
                print(self.scanned)
                self.status_changed()

        self.start(loop=async_loop)

    def reader_handle_result(self, result):
        self.log("Submitted code: {}".format(result))
        self.currentValue = result
        self.status_changed()

    def reader_status_changed(self, result):
        self.log("Current status: {}".format(result))
        self.currentValue = result
        self.status_changed()


scan_codes = {
    # Scancode: ASCIICode
    0: None, 1: u'ESC', 2: u'1', 3: u'2', 4: u'3', 5: u'4', 6: u'5', 7: u'6', 8: u'7', 9: u'8',
    10: u'9', 11: u'0', 12: u'-', 13: u'=', 14: u'BKSP', 15: u'TAB', 16: u'q', 17: u'w', 18: u'e', 19: u'r',
    20: u't', 21: u'y', 22: u'u', 23: u'i', 24: u'o', 25: u'p', 26: u'[', 27: u']', 28: u'CRLF', 29: u'LCTRL',
    30: u'a', 31: u's', 32: u'd', 33: u'f', 34: u'g', 35: u'h', 36: u'j', 37: u'k', 38: u'l', 39: u';',
    40: u'"', 41: u'`', 42: u'LSHFT', 43: u'\\', 44: u'z', 45: u'x', 46: u'c', 47: u'v', 48: u'b', 49: u'n',
    50: u'm', 51: u',', 52: u'.', 53: u'/', 54: u'RSHFT', 56: u'LALT', 57: u' ', 100: u'RALT'
}

caps_codes = {
    0: None, 1: u'ESC', 2: u'!', 3: u'@', 4: u'#', 5: u'$', 6: u'%', 7: u'^', 8: u'&', 9: u'*',
    10: u'(', 11: u')', 12: u'_', 13: u'+', 14: u'BKSP', 15: u'TAB', 16: u'Q', 17: u'W', 18: u'E', 19: u'R',
    20: u'T', 21: u'Y', 22: u'U', 23: u'I', 24: u'O', 25: u'P', 26: u'{', 27: u'}', 28: u'CRLF', 29: u'LCTRL',
    30: u'A', 31: u'S', 32: u'D', 33: u'F', 34: u'G', 35: u'H', 36: u'J', 37: u'K', 38: u'L', 39: u':',
    40: u'\'', 41: u'~', 42: u'LSHFT', 43: u'|', 44: u'Z', 45: u'X', 46: u'C', 47: u'V', 48: u'B', 49: u'N',
    50: u'M', 51: u'<', 52: u'>', 53: u'?', 54: u'RSHFT', 56: u'LALT', 57: u' ', 100: u'RALT'
}


def read_barcode(device_path):
    dev = InputDevice(device_path)
    dev.grab()  # grab provides exclusive access to the device

    x = ''
    caps = False

    for event in dev.read_loop():
        if event.type == ecodes.EV_KEY:
            data = categorize(event)  # Save the event temporarily to introspect it
            if data.scancode == 42:
                if data.keystate == 1:
                    caps = True
                if data.keystate == 0:
                    caps = False

            if data.keystate == 1:  # Down events only
                if caps:
                    key_lookup = u'{}'.format(caps_codes.get(data.scancode)) or u'UNKNOWN:[{}]'.format(
                        data.scancode)  # Lookup or return UNKNOWN:XX
                else:
                    key_lookup = u'{}'.format(scan_codes.get(data.scancode)) or u'UNKNOWN:[{}]'.format(
                        data.scancode)  # Lookup or return UNKNOWN:XX

                if (data.scancode != 42) and (data.scancode != 28):
                    x += key_lookup

                if (data.scancode == 28):
                    return x


if __name__ == "__main__":
    device = Scanner()
    device.main()
