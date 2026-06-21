from sqlalchemy import (
    CheckConstraint,
    Column,
    Integer,
    String,
    ForeignKey,
    DateTime,
)
from sqlalchemy.sql import func

from app.models.base import Base


class WatchEvent(Base):
    __tablename__ = "watch_events"

    __table_args__ = (
        CheckConstraint(
            "position_seconds >= 0",
            name="positive_position_seconds",
        ),
    )

    id = Column(
        Integer,
        primary_key=True,
    )

    user_id = Column(
        Integer,
        ForeignKey("users.id"),
        nullable=False,
    )

    movie_id = Column(
        Integer,
        ForeignKey("movies.id"),
        nullable=False,
    )

    event_type = Column(
        String,
        nullable=False,
    )

    position_seconds = Column(
        Integer,
        nullable=False,
    )

    created_at = Column(
        DateTime(timezone=True),
        server_default=func.now(),
        nullable=False,
    )
