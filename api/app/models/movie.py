from sqlalchemy import String, DateTime, func
from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column
from app.config import settings
from app.models.base import Base


class Movie(Base):
    __tablename__ = "movies"

    id: Mapped[int] = mapped_column(primary_key=True)

    title: Mapped[str] = mapped_column(String(255))

    description: Mapped[str | None]

    filename: Mapped[str] = mapped_column(String(255))

    created_at: Mapped[str] = mapped_column(
        DateTime(timezone=True),
        server_default=func.now(),
    )

    poster_filename: Mapped[str] = mapped_column(String(255), nullable=True)

    @property
    def poster_url(self):
        if not self.poster_filename:
            return None
        return f"{settings.STREAM_BASE_URL}/posters/{self.poster_filename}"
