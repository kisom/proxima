# UPS power pack

import os
import time
import smbus2
import logging
from ina219 import INA219,DeviceRangeError

# Define I2C bus
DEVICE_BUS = 1

# Define device i2c slave address.
DEVICE_ADDR = 0x17

# Set the threshold of UPS automatic power-off to prevent damage caused by battery over-discharge, unit: mV.
PROTECT_VOLT = 3700  

# Set the sample period, Unit: min default: 2 min.
SAMPLE_TIME = 2

class UPS:
    def __init__(self):
        self.v_b = None
        self.v_w = None
        self.c_b = None
        self.c_w = None
        self.p_b = None
        self.p_w = None
        self.m = None
        self.batt = INA219(0.005, busnum=DEVICE_BUS, address=0x45)
        self.wall = INA219(0.00725, busnum=DEVICE_BUS, address=0x40)
        self.bus = smbus2.SMBus(DEVICE_BUS)

        self.s_b = None
        self.s_w = None

        self.batt.configure()
        self.wall.configure()
        self.refresh()

    def _collect(self, source):
        v = source.voltage()
        c = source.current()
        p = source.power()
        return (v, c, p)

    def refresh(self):
        try:
            (self.v_b, self.c_b, self.p_b) = self._collect(self.batt)
            (self.v_w, self.c_w, self.p_w) = self._collect(self.wall)
            self._state()
            self.m = None
        except Exception as exc:
            self.m = str(exc)

    def _state(self):
        if self.c_b > 0:
            self.s_b = 'C'
        else:
            self.s_b = 'D'
        self._charge_state()

    def _charge_state(self):
        buf = []
        buf.append(0x00)
        for i in range(1, 255):
            buf.append(self.bus.read_byte_data(DEVICE_ADDR, i))
        if (buf[8] << 8 | buf[7]) > 4000:
            self.s_w = 'USBC'
        elif (buf[10] << 8 | buf[9])> 4000:
            self.s_w = 'USBM'
        else:
            self.s_w = "----"

        if ((self.v_b * 1000) < (PROTECT_VOLT)):
            self.m = 'Battery almost dead!'

    def has_error(self):
        return self.m is not None

    def clear_error(self):
        self.m = None

    def __str__(self):
        if self.m:
            return f"!!!: {self.m}"
        return f"""BATT: V: {self.v_b:0.1f}V C: {self.v_b:0.1f}mA P: {self.p_b:0.1f}mW
WALL: V: {self.v_w:0.1f}V C: {self.v_w:0.1f}mA P: {self.p_w:0.1f}mW
STATE: B:{self.s_b} S:{self.s_w}"""

    def lines(self):
        return str(self).split('\n')
