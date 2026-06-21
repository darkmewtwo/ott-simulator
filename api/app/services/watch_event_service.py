from app.repositories.watch_event_repository import (
    WatchEventRepository,
)

from app.schemas.watch_event import (
    WatchEventCreate,
    WatchEventResponse,
)

from app.models.user import User


class WatchEventService:
    def __init__(
        self,
        repo: WatchEventRepository,
    ):
        self.repo = repo

    def create_event(
        self,
        payload: WatchEventCreate,
        user: User,
    ) -> WatchEventResponse:

        event = self.repo.create(
            user_id=user.id,
            movie_id=payload.movie_id,
            event_type=payload.event_type,
            position_seconds=payload.position_seconds,
        )

        return WatchEventResponse.model_validate(event)
