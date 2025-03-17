from enum import StrEnum, auto


class HeicMode(StrEnum):
    CONVERT = auto()
    COPY = auto()

    @staticmethod
    def list():
        return list(map(lambda c: c.value, HeicMode))
