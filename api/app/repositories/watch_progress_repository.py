from sqlalchemy.orm import Session

from app.models.watch_progress import (
    WatchProgress,
)


class WatchProgressRepository:
    def __init__(
        self,
        db: Session,
    ):
        self.db = db

    def get(
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

    def upsert(
        self,
        user_id: int,
        movie_id: int,
        position_seconds: int,
    ):
        progress = self.get(
            user_id,
            movie_id,
        )

        if progress:
            progress.last_position_seconds = position_seconds

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
