from datetime import date
from sqlalchemy import String, DateTime, func, Integer, text, Date, JSON
from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column
from app.constants.movie import MovieStatus
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

    duration_seconds: Mapped[int] = mapped_column(
        Integer, default=0, server_default=text("0")
    )

    status: Mapped[MovieStatus] = mapped_column(
        String(20),
        nullable=False,
        default="PENDING",
        server_default="PENDING",
    )

    hls_playlist_path: Mapped[str | None] = mapped_column(
        String,
        nullable=True,
    )
    # Metadata

    release_date: Mapped[date | None] = mapped_column(
        Date,
        nullable=True,
    )

    language: Mapped[str | None] = mapped_column(
        String(50),
        nullable=True,
    )

    genres: Mapped[list[str] | None] = mapped_column(
        JSON,
        nullable=True,
    )

    age_rating: Mapped[str | None] = mapped_column(
        String(20),
        nullable=True,
    )

    director: Mapped[str | None] = mapped_column(
        String(255),
        nullable=True,
    )

    cast: Mapped[list[str] | None] = mapped_column(
        JSON,
        nullable=True,
    )

    @property
    def poster_url(self):
        if not self.poster_filename:
            return None
        return f"/posters/{self.poster_filename}"
