from datetime import datetime, timedelta, timezone

from jose import jwt

from app.config import settings


def create_stream_token(
    user_id: int,
    movie_id: int,
) -> str:
    expire = datetime.now(timezone.utc) + timedelta(
        hours=settings.STREAM_TOKEN_EXPIRE_HOURS,
    )

    payload = {
        "user_id": user_id,
        "movie_id": movie_id,
        "exp": expire,
    }

    return jwt.encode(
        payload,
        settings.STREAM_SECRET_KEY,
        algorithm=settings.ALGORITHM,
    )
