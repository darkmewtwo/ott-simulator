from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

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
def get_movie(movie_id: int, service: MovieService = Depends(get_service)):
    return service.get_movie(movie_id)


@router.post("", response_model=MovieResponse)
def create_movie(payload: MovieCreate, service: MovieService = Depends(get_service)):
    return service.create_movie(payload)