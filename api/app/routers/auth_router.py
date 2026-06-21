from fastapi import APIRouter, Depends

from app.dependencies import get_db

from sqlalchemy.orm import Session

from app.repositories.user_repository import UserRepository

from app.services.auth_service import AuthService

from app.schemas.user import TokenResponse, UserCreate, UserLogin, UserResponse

from app.core.dependencies import (
    get_current_user,
)

from app.models.user import User

router = APIRouter(prefix="/auth", tags=["auth"])


def get_service(db: Session = Depends(get_db)) -> AuthService:

    user_repo = UserRepository(db)

    return AuthService(user_repo)


@router.post("/register", response_model=UserResponse)
def register(payload: UserCreate, service: AuthService = Depends(get_service)):
    return service.register(payload)


@router.post(
    "/login",
    response_model=TokenResponse,
)
def login(
    payload: UserLogin,
    service: AuthService = Depends(get_service),
):
    return service.login(payload)


@router.get(
    "/me",
    response_model=UserResponse,
)
def me(
    current_user: User = Depends(get_current_user),
):
    return current_user
