from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    DATABASE_URL: str = "postgresql+psycopg://ott:ott@postgres:5432/ott"
    STREAM_BASE_URL: str = "http://localhost:8180"
    STREAM_SECRET_KEY: str = "stream-secret-key-change-later"
    ALGORITHM: str = "HS256"

    STREAM_TOKEN_EXPIRE_HOURS: int = 4


settings = Settings()
