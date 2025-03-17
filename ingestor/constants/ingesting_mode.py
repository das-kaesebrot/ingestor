from enum import StrEnum, auto


class IngestingMode(StrEnum):
    MOVE = auto()
    COPY = auto()

    @staticmethod
    def list():
        return list(map(lambda c: c.value, IngestingMode))
