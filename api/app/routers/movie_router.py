from datetime import date

from fastapi import APIRouter, Depends, Form, File, Response, UploadFile
from sqlalchemy.orm import Session
from typing import Optional

from app.core.playback import create_stream_token
from app.core.dependencies import get_current_user
from app.models.user import User
from app.dependencies import get_db
from app.repositories.movie_repository import MovieRepository
from app.services.movie_service import MovieService
from app.schemas.movie import MovieCreate, MovieDetailsResponse, MovieResponse

router = APIRouter(prefix="/movies", tags=["movies"])


def get_service(db: Session = Depends(get_db)):
    repo = MovieRepository(db)
    return MovieService(repo)


@router.get("", response_model=list[MovieResponse])
def list_movies(service: MovieService = Depends(get_service)):
    return service.list_movies()


@router.get("/{movie_id}", response_model=MovieDetailsResponse)
def get_movie(
    movie_id: int,
    response: Response,
    current_user: User = Depends(get_current_user),
    service: MovieService = Depends(get_service),
):
    stream_token = create_stream_token(
        user_id=current_user.id,
        movie_id=movie_id,
    )

    response.set_cookie(
        key="stream_token",
        value=stream_token,
        httponly=True,
        secure=False,  # localhost
        samesite="lax",
        max_age=4 * 60 * 60,
    )

    return service.get_movie(movie_id)


@router.post("", response_model=MovieResponse)
def create_movie(payload: MovieCreate, service: MovieService = Depends(get_service)):
    return service.create_movie(payload)


@router.post("/upload", response_model=MovieResponse)
def upload_movie(
    title: str = Form(...),
    description: str = Form(""),
    release_date: date | None = Form(None),
    language: str | None = Form(None),
    genres: list | None = Form(None),      # JSON array
    age_rating: str | None = Form(None),
    director: str | None = Form(None),
    cast: list | None = Form(None),        # JSON array
    video_file: UploadFile = File(...),
    poster_file: Optional[UploadFile] = File(None),
    service: MovieService = Depends(get_service),
):
    return service.upload_movie(
        title=title,
        description=description,
        release_date=release_date,
        language=language,
        genres=genres,
        age_rating=age_rating,
        director=director,
        cast=cast,
        file=video_file,
        poster_file=poster_file,
    )
