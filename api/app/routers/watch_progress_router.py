from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from app.core.dependencies import get_current_user

from app.dependencies import get_db

from app.repositories.watch_progress_repository import (
    WatchProgressRepository,
)

from app.schemas.watch_progress import WatchProgressResponse

from app.services.watch_progress_service import WatchProgressService

from app.models.user import User

router = APIRouter(
    prefix="/watch_progress",
    tags=["watch_progress"],
)


def get_service(
    db: Session = Depends(get_db),
):
    repo = WatchProgressRepository(db)
    return WatchProgressService(repo)


@router.get("/progress", response_model=list[WatchProgressResponse])
def get_progress(
    current_user: User = Depends(get_current_user),
    service: WatchProgressService = Depends(get_service),
):
    return service.get_user_progress(current_user)


@router.get("/{movie_id}", response_model=WatchProgressResponse | None)
def get_movie_progress(
    movie_id: int,
    current_user: User = Depends(get_current_user),
    service: WatchProgressService = Depends(get_service),
):
    return service.get_movie_progress(
        current_user.id,
        movie_id,
    )
