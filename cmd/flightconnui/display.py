import platform
import time
from ups import UPS

SHIP_NAME = "KISS Beauty of the Deep Night"

# SW1 = 16 # button 4: buggy
# SW2 = 26 # button 3: buggy
SW3 = 20 # button 2
SW4 = 21 # button 1
SW = [SW4, SW3]

BATT_PAGE = 2

def is_raspberry_pi() -> bool:
    return platform.machine() in ('armv7l', 'armv6l')

if is_raspberry_pi():
    from papirus import Papirus
    from papirus import PapirusTextPos
    import RPi.GPIO as GPIO

class ConsoleDisplay:
    def __init__(self):
        pass

    def display(self, lines):
        print(SHIP_NAME)
        for line in lines:
            print(line)

    def min_interval(self):
        return 1

    def check_buttons(self):
        pass

    def show_battery_status(self):
        pass


class PapirusDisplay:
    max_page = 2

    def __init__(self, rot=0):
        self.text = PapirusTextPos(autoUpdate=False, rotation=rot)
        self.text.papirus.fast_update()
        self.page = 0
        self.cached = None
        self.scroll = True
        self.pressed = None
        self.update = True
        self.pause = False
        self.ups = UPS()

        GPIO.setmode(GPIO.BCM)
        for sw in SW:
            GPIO.setup(sw, GPIO.IN)

    def display(self, m):
        if not m or self.pause:
            return
        if self.pause:
            self.pause = True
        self.text.Clear()
        linepos = 0
        self.text.AddText(SHIP_NAME + f"[{self._status_char()}]", 0, 0, size=14)
        if self.page == BATT_PAGE:
            self.ups.refresh()
            lines = self.ups.lines()
            lines.insert(0, f"Ship time: {m.clock.o}")
        else:
            lines = m.lines(self.page)
        for line in lines:
            parts = line.split(':', 1)
            x = 0
            for part in parts:
                linepos += 12
                self.text.AddText(part, x, linepos, size=14)
                x = 10
        self.text.WriteAll()
        self.cached = m
        if not self.update:
            self.pause = True
        if self.scroll:
            self.page_forward()

    def page_forward(self):
        if not self.scroll:
            return
        self.page += 1
        self._adj_page()

    def page_back(self):
        if not self.scroll:
            return
        self.page -= 1
        self._adj_page()

    def _adj_page(self):
        if self.page < 0:
            self.page = self.max_page
        elif self.page > self.max_page:
            self.page = 0

    def min_interval(self):
        return 60

    def _bounce_ok(self, display=False):
        if self.pressed is not None:
            now = time.time()
            delta = (now - self.pressed)
            if display:
                print(f'pressed: {self.pressed:0.1f}')
                print(f'now: {now:0.1f}')
                print(f'delta: {delta:0.3f}')
            if delta < 1.0:
                return False
        return True

    def _clear_button_press(self):
        if self._bounce_ok():
            self.pressed = None

    def check_buttons(self):
        buttons = []
        for i in range(len(SW)):
            if not GPIO.input(SW[i]):
                buttons.append((i, SW[i]))
        if buttons:
            button = min(buttons, key = lambda x: x[1])
            return button[0]

    def _status_char(self):
        if not self.update:
            return '#'
        if not self.scroll:
            return '*'
        return ' '

    def handle_buttons(self, display=False):
        button = self.check_buttons()
        # buttons 3 and 4 are buggy
        if button is None:
            self._clear_button_press()
            return
        if not self._bounce_ok(display):
            self.pressed = time.time()
            return
        if button == 1:
            self.update = not self.update
            self.pause = False
        if button == 0:
            self.scroll = not self.scroll
        self.display(self.cached)
        self.pressed = time.time()
        

def get_display():
    if is_raspberry_pi():
        return PapirusDisplay()
    return ConsoleDisplay()

