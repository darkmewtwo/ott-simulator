from app.schemas.watch_progress import (
    ContinueWatchingResponse,
    WatchHistoryMovieResponse,
    WatchHistoryProgressResponse,
    WatchHistoryResponse,
)
from app.models.user import User

from app.repositories.watch_progress_repository import WatchProgressRepository


class WatchProgressService:
    def __init__(
        self,
        repo: WatchProgressRepository,
    ):
        self.repo = repo

    def get_user_progress(
        self,
        user: User,
    ):
        return self.repo.get_by_user(user.id)

    def get_movie_progress(
        self,
        user_id: int,
        movie_id: int,
    ):
        return self.repo.get_by_user_and_movie(
            user_id,
            movie_id,
        )

    def get_continue_watching(
        self,
        user_id: int,
    ):
        rows = self.repo.get_continue_watching(user_id)

        return [
            ContinueWatchingResponse(
                movie_id=movie.id,
                title=movie.title,
                poster_url=movie.poster_url,
                last_position_seconds=progress.last_position_seconds,
            )
            for progress, movie in rows
        ]

    def get_watch_history(
        self,
        user: User,
    ):
        rows = self.repo.get_watch_history(user.id)

        return [
            WatchHistoryResponse(
                movie=WatchHistoryMovieResponse(
                    id=movie.id,
                    title=movie.title,
                    poster_url=movie.poster_url,
                    duration_seconds=movie.duration_seconds,
                ),
                progress=WatchHistoryProgressResponse(
                    last_position_seconds=progress.last_position_seconds,
                    watch_percentage=(
                        round(
                            progress.last_position_seconds
                            / movie.duration_seconds
                            * 100,
                            2,
                        )
                        if movie.duration_seconds > 0
                        else 0
                    ),
                    started_at=progress.started_at,
                    last_watched_at=progress.last_watched_at,
                    is_completed=progress.is_completed,
                ),
            )
            for progress, movie in rows
        ]
