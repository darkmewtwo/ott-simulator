from app.schemas.watch_progress import ContinueWatchingResponse
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
        print(user_id)
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
