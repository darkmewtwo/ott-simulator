from pathlib import Path
import shutil
import uuid

from fastapi import HTTPException, UploadFile

from app.config import settings
from app.models.movie import Movie
from app.repositories.movie_repository import MovieRepository
from app.schemas.movie import MovieCreate


MEDIA_DIR = Path("/media")


class MovieService:
    def __init__(self, repo: MovieRepository):
        self.repo = repo

    def list_movies(self):
        return self.repo.list_movies()

    def get_movie(self, movie_id: int):
        movie = self.repo.get_movie(movie_id)
        movie.stream_url = f"{settings.STREAM_BASE_URL}/stream/{movie.filename}"
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

    def upload_movie(
        self,
        title: str,
        description: str,
        file: UploadFile,
    ):

        if not file.filename:
            raise HTTPException(
                status_code=400,
                detail="No file provided",
            )

        extension = Path(file.filename).suffix

        filename = f"{uuid.uuid4()}{extension}"

        MEDIA_DIR.mkdir(exist_ok=True)

        filepath = MEDIA_DIR / filename

        with open(filepath, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)

        movie = Movie(
            title=title,
            description=description,
            filename=filename,
        )

        return self.repo.create_movie(movie)
