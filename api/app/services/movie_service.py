from fastapi import HTTPException

from app.models.movie import Movie
from app.repositories.movie_repository import MovieRepository
from app.schemas.movie import MovieCreate


class MovieService:

    def __init__(self, repo: MovieRepository):
        self.repo = repo

    def list_movies(self):
        return self.repo.list_movies()

    def get_movie(self, movie_id: int):
        movie = self.repo.get_movie(movie_id)
        if not movie:
            raise HTTPException(status_code=404, detail="Movie not found")
        return movie

    def create_movie(self, payload: MovieCreate):
        movie = Movie(
            title=payload.title,
            description=payload.description,
            filename=payload.filename,
        )
        return self.repo.create_movie(movie)