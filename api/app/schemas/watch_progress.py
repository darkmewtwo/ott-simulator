from datetime import datetime

from pydantic import BaseModel, ConfigDict


class WatchProgressResponse(BaseModel):
    movie_id: int
    last_position_seconds: int
    is_completed: bool

    model_config = ConfigDict(from_attributes=True)


class ContinueWatchingResponse(BaseModel):
    movie_id: int
    title: str
    poster_url: str | None
    last_position_seconds: int

    model_config = ConfigDict(from_attributes=True)


class WatchHistoryMovieResponse(BaseModel):
    id: int
    title: str
    poster_url: str | None
    duration_seconds: int

    model_config = ConfigDict(from_attributes=True)


class WatchHistoryProgressResponse(BaseModel):
    last_position_seconds: int
    watch_percentage: float
    started_at: datetime | None
    last_watched_at: datetime | None
    is_completed: bool

    model_config = ConfigDict(from_attributes=True)


class WatchHistoryResponse(BaseModel):
    movie: WatchHistoryMovieResponse
    progress: WatchHistoryProgressResponse
