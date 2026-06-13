from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    DATABASE_URL: str = (
        "postgresql+psycopg://ott:ott@postgres:5432/ott"
    )


settings = Settings()