from sqlalchemy.orm import Session

from app.constants.watch_event import WatchEventType
from app.models.watch_event import WatchEvent


class WatchEventRepository:
    def __init__(
        self,
        db: Session,
    ):
        self.db = db

    def create(
        self,
        user_id: int,
        movie_id: int,
        event_type: WatchEventType,
        position_seconds: int,
    ) -> WatchEvent:

        event = WatchEvent(
            user_id=user_id,
            movie_id=movie_id,
            event_type=event_type.value,
            position_seconds=position_seconds,
        )

        self.db.add(event)
        self.db.commit()
        self.db.refresh(event)

        return event
