from datetime import datetime, timezone
from sqlalchemy.orm import Session

from app.models.movie import Movie
from app.models.watch_progress import (
    WatchProgress,
)


class WatchProgressRepository:
    def __init__(
        self,
        db: Session,
    ):
        self.db = db

    def upsert(
        self,
        user_id: int,
        movie_id: int,
        position_seconds: int,
        is_completed: bool = False,
        update_watch_time: bool = False,
    ):
        progress = self.get_by_user_and_movie(
            user_id,
            movie_id,
        )
        now = datetime.now(timezone.utc)
        if progress:
            progress.is_completed = is_completed
            progress.last_watched_at = now

            if update_watch_time:
                previous_position = progress.last_position_seconds
                delta = position_seconds - previous_position

                if delta < 0:
                    delta = 0

                if delta > 15:
                    delta = 0

                progress.total_watch_time_seconds += delta

            progress.last_position_seconds = position_seconds

        else:
            progress = WatchProgress(
                user_id=user_id,
                movie_id=movie_id,
                last_position_seconds=position_seconds,
                started_at=now,
                last_watched_at=now,
            )

            self.db.add(progress)

        self.db.commit()

        self.db.refresh(progress)

        return progress

    def get_by_user(
        self,
        user_id: int,
    ):
        return (
            self.db.query(WatchProgress).filter(WatchProgress.user_id == user_id).all()
        )

    def get_by_user_and_movie(
        self,
        user_id: int,
        movie_id: int,
    ):
        return (
            self.db.query(WatchProgress)
            .filter(
                WatchProgress.user_id == user_id,
                WatchProgress.movie_id == movie_id,
            )
            .first()
        )

    def get_continue_watching(
        self,
        user_id: int,
        limit: int = 20,
    ):
        return (
            self.db.query(
                WatchProgress,
                Movie,
            )
            .join(
                Movie,
                Movie.id == WatchProgress.movie_id,
            )
            .filter(
                WatchProgress.user_id == user_id,
                WatchProgress.is_completed == False,  # noqa
            )
            .order_by(WatchProgress.updated_at.desc())
            .limit(limit)
            .all()
        )

    def get_watch_history(
        self,
        user_id: int,
    ) -> list[WatchProgress]:
        return (
            self.db.query(
                WatchProgress,
                Movie,
            )
            .join(
                Movie,
                Movie.id == WatchProgress.movie_id,
            )
            .filter(
                WatchProgress.user_id == user_id,
            )
            .order_by(
                WatchProgress.last_watched_at.desc(),
            )
            .all()
        )
