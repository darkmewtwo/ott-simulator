from datetime import datetime
from pydantic import BaseModel, field_serializer


class MovieCreate(BaseModel):
    title: str
    description: str | None = None
    filename: str


class MovieResponse(BaseModel):
    id: int
    title: str
    description: str | None
    filename: str
    created_at: datetime

    @field_serializer("created_at")
    def serialize_created_at(self, value: datetime):
        return value.isoformat()

    class Config:
        from_attributes = True