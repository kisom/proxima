#!/usr/bin/env python3

import json
import zmq

import display
import physics

ctx = zmq.Context().instance()

def getkey(m, k, d=None):
    if k in m:
        return m[k]
    return d

class State:
    def __init__(self, s):
        self.x = int(getkey(s, "x", 0))
        self.v = int(getkey(s, "v", 0))

class Clock:
    def __init__(self, c):
        self.l = getkey(c, "launched")
        self.oet = int(getkey(c, "observer_et", 0))
        self.ret = int(getkey(c, "relative_et", 0))
        self.o = getkey(c, "observer")
        self.r = getkey(c, "relative")

class Mission:
    max_page = 1

    def __init__(self, s, c, a):
        self.action = a
        self.state = s
        self.clock = c

    def display(self, page):
        s = f"""Phase: {self.action}
 Ship time: {self.clock.o}
Earth time: {self.clock.r}
Velocity: {self.state.v/1000:0.1f} km/s ({physics.v_to_c(self.state.v):0.3f}c)
"""
        if page == 0:
            s += f"""Distance traveled: {physics.distance_string(self.state.x)}
Distance remaining: {physics.distance_string(physics.proxima-self.state.x)}"""
        elif page == 1:
            s += f"""Ship elapsed: {physics.time_string(self.clock.oet)}
Relative elapsed: {physics.time_string(self.clock.ret)}"""
        return s

    def lines(self, page):
        return self.display(page).split("\n")


def from_json(msg):
    update = json.loads(msg)
    return Mission(
        State(update["state"]),
        Clock(update["clock"]),
        update["action"] 
    )

def recv_update(sock):
    return from_json(sock.recv())
        
def connect(addr="tcp://localhost:4247", topic=""):
    sock = ctx.socket(zmq.SUB)
    sock.setsockopt(zmq.CONFLATE, 1)
    sock.connect(addr)
    sock.subscribe(topic)
    return sock

def _conn_and_rx():
    return recv_update(connect())


def main():
    sock = connect()
    disp = display.get_display()
    while True:
        m = recv_update(sock)
        disp.display(m)
        delta = max(physics.update_time(m), disp.min_interval())
        physics.sleep(delta, action=disp.handle_buttons)

if __name__ == '__main__':
    main()
