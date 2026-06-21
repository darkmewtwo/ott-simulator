from fastapi import Depends
from fastapi import HTTPException
from fastapi.security import HTTPBearer
from fastapi.security import HTTPAuthorizationCredentials

from sqlalchemy.orm import Session

from app.dependencies import get_db

from app.repositories.user_repository import (
    UserRepository,
)

from app.core.security import (
    decode_access_token,
)

security = HTTPBearer()


def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(security),
    db: Session = Depends(get_db),
):
    token = credentials.credentials

    payload = decode_access_token(token)

    if not payload:
        raise HTTPException(
            status_code=401,
            detail="Invalid token",
        )

    user_id = int(payload["sub"])

    user_repo = UserRepository(db)

    user = user_repo.get_by_id(user_id)

    if not user:
        raise HTTPException(
            status_code=401,
            detail="User not found",
        )

    return user
