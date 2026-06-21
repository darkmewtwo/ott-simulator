from pydantic import BaseModel

from app.constants.watch_event import WatchEventType


class WatchEventCreate(BaseModel):
    movie_id: int
    event_type: WatchEventType
    position_seconds: int


class WatchEventResponse(BaseModel):
    id: int
    user_id: int
    movie_id: int
    event_type: WatchEventType
    position_seconds: int

    model_config = {
        "from_attributes": True,
    }
