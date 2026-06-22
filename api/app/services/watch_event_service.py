from app.repositories.watch_event_repository import (
    WatchEventRepository,
)
from app.repositories.watch_progress_repository import WatchProgressRepository

from app.schemas.watch_event import (
    WatchEventCreate,
    WatchEventResponse,
)

from app.models.user import User


class WatchEventService:
    def __init__(
        self,
        repo: WatchEventRepository,
        progress_repo: WatchProgressRepository,
    ):
        self.repo = repo
        self.progress_repo = progress_repo

    def create_event(
        self,
        payload: WatchEventCreate,
        user: User,
    ) -> WatchEventResponse:

        if payload.event_type.value in (
            "PAUSE",
            "SEEK",
            "STOP",
            "COMPLETE",
        ):
            self.progress_repo.upsert(
                user_id=user.id,
                movie_id=payload.movie_id,
                position_seconds=payload.position_seconds,
            )
        event = self.repo.create(
            user_id=user.id,
            movie_id=payload.movie_id,
            event_type=payload.event_type,
            position_seconds=payload.position_seconds,
        )

        return WatchEventResponse.model_validate(event)

    def get_user_progress(
        self,
        user: User,
    ):
        return self.progress_repo.get_by_user(user.id)
