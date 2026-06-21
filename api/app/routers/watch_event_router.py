from fastapi import (
    APIRouter,
    Depends,
)

from sqlalchemy.orm import Session

from app.dependencies import get_db

from app.models.user import User

from app.core.dependencies import (
    get_current_user,
)

from app.repositories.watch_event_repository import (
    WatchEventRepository,
)

from app.services.watch_event_service import (
    WatchEventService,
)

from app.schemas.watch_event import (
    WatchEventCreate,
    WatchEventResponse,
)

router = APIRouter(
    prefix="/events",
    tags=["events"],
)


def get_service(
    db: Session = Depends(get_db),
):
    repo = WatchEventRepository(db)

    return WatchEventService(repo)


@router.post(
    "",
    response_model=WatchEventResponse,
)
def create_event(
    payload: WatchEventCreate,
    current_user: User = Depends(get_current_user),
    service: WatchEventService = Depends(get_service),
):
    return service.create_event(
        payload,
        current_user,
    )
