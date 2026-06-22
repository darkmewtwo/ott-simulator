from pydantic import BaseModel, ConfigDict


class WatchProgressResponse(BaseModel):
    movie_id: int
    last_position_seconds: int

    model_config = ConfigDict(from_attributes=True)
