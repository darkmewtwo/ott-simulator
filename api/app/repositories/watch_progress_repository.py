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
        is_completed:bool = False
    ):
        progress = self.get_by_user_and_movie(
            user_id,
            movie_id,
        )

        if progress:
            progress.last_position_seconds = position_seconds
            progress.is_completed = is_completed

        else:
            progress = WatchProgress(
                user_id=user_id,
                movie_id=movie_id,
                last_position_seconds=position_seconds,
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
                Movie.id
                == WatchProgress.movie_id,
            )
            .filter(
                WatchProgress.user_id
                == user_id,
                WatchProgress.is_completed == False # noqa
                ,
            )
            .order_by(
                WatchProgress.updated_at.desc()
            )
            .limit(limit)
            .all()
        )