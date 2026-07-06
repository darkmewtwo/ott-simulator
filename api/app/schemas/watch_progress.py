from pydantic import BaseModel, ConfigDict


class WatchProgressResponse(BaseModel):
    movie_id: int
    last_position_seconds: int

    model_config = ConfigDict(from_attributes=True)


class ContinueWatchingResponse(BaseModel):
    movie_id: int
    title: str
    poster_url: str | None
    last_position_seconds: int

    model_config = ConfigDict(from_attributes=True)
