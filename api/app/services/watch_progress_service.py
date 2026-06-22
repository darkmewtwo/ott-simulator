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
