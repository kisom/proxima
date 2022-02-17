import math
import time

c = 299792458
lym = 9.46073047258e+15
au = 149597870691

# convert lightyears to meters
def ly_to_m(ly):
    return ly * lym

# convert meters to lightyears
def m_to_ly(m):
    return m / lym

# return velocity as a percent of lightspeed
def v_to_c(v):
    return v / c

# return a velocity from a percent of lightspeed
def c_to_v(pct):
    return pct * c

# convert meters to au
def m_to_au(m):
    return m / au

# convert au to meters
def au_to_m(x):
    return x * au


proxima      = ly_to_m(4.247)
hundredth_ly = ly_to_m(0.01)
tenth_au     = au_to_m(0.1)
hundredth_au = au_to_m(0.01)
jupiter      = au_to_m(5.0)
fifty_au     = au_to_m(50.0)
heliopause   = au_to_m(960.78)

def distance_string(x):
    if x > hundredth_ly:
        return f"{m_to_ly(x):0.4f} ly"
    # technically, I think au should be used in relation to the sun...
    if x > tenth_au:
        return f"{m_to_au(x):0.1f} au"
    if x > hundredth_au:
        return f"{m_to_au(x):0.2f} au"
    return f"{x/1000} km"

sec_to_day = 86400
sec_to_year = 31557600 

def time_string(x):
    x = int(x)
    s = []
    years = x // sec_to_year
    if years > 0:
        s.append(f"{int(years)}y")
        x -= (math.floor(years) * sec_to_year)
    days = x // sec_to_day
    if days > 0:
        s.append(f"{int(days)}d")
        x -= (math.floor(days) * sec_to_day)
    if years < 1 and days < 7:
        hours = x // 3600
        if hours > 0:
            s.append(f"{int(hours)}h")
            x -= (math.floor(hours) * 3600)
        mins = x // 60
        if mins > 0:
            s.append(f"{int(mins)}m")
            x -= (math.floor(mins) * 60)
        if x > 0:
            s.append(f"{int(x)}s")

    return ' '.join(s)

def update_time(m):
    if m.action == "accelerating":
        if m.state.x < jupiter:
            return 1
        elif m.state.x < heliopause:
            return 60
    elif m.action == "decelerating":
        remaining = proxima - m.state.x
        if remaining < fifty_au:
            return 60
        if remaining < jupiter:
            return 1
    return 3600

def sleep(n, action=None):
    start = time.time()
    stop = start + n
    while time.time() < stop:
        if action:
            action()
        time.sleep(1)

