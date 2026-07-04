from pathlib import Path
import shutil
import uuid

from fastapi import HTTPException, UploadFile

from app.constants.movie import MovieStatus
from app.config import settings
from app.models.movie import Movie
from app.repositories.movie_repository import MovieRepository
from app.schemas.movie import MovieCreate


MEDIA_DIR = Path("/media/movies")
POSTER_DIR = Path("/media/posters")


class MovieService:
    def __init__(self, repo: MovieRepository):
        self.repo = repo

    def list_movies(self):
        return self.repo.list_movies()

    def get_movie(self, movie_id: int):
        movie = self.repo.get_movie(movie_id)
        if not movie:
            raise HTTPException(status_code=404, detail="Movie not found")
        
        movie.stream_url = f"{settings.STREAM_BASE_URL}/stream/{movie.filename}"
        movie.hls_url = None
            
        if movie.status == MovieStatus.READY:
            movie.hls_url = f"{settings.STREAM_BASE_URL}/hls/{movie.id}/index.m3u8"
        
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
        poster_file: UploadFile | None = None,
    ):
        video_file = file
        if not video_file.filename:
            raise HTTPException(
                status_code=400,
                detail="No file provided",
            )

        video_extension = Path(video_file.filename).suffix

        video_filename = f"{uuid.uuid4()}{video_extension}"

        MEDIA_DIR.mkdir(exist_ok=True)

        filepath = MEDIA_DIR / video_filename

        with open(filepath, "wb") as buffer:
            shutil.copyfileobj(video_file.file, buffer)

        poster_filename = None

        if poster_file and poster_file.filename:
            poster_extension = Path(poster_file.filename).suffix

            poster_filename = f"{uuid.uuid4()}{poster_extension}"

            poster_filepath = POSTER_DIR / poster_filename

            with open(poster_filepath, "wb") as buffer:
                shutil.copyfileobj(poster_file.file, buffer)

        movie = Movie(
            title=title,
            description=description,
            filename=video_filename,
            poster_filename=poster_filename,
        )

        return self.repo.create_movie(movie)
