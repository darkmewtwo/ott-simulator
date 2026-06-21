from enum import Enum


class WatchEventType(str, Enum):
    PLAY = "PLAY"
    PAUSE = "PAUSE"
    SEEK = "SEEK"
    STOP = "STOP"
    COMPLETE = "COMPLETE"
