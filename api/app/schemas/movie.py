from datetime import datetime
from pydantic import BaseModel, field_serializer, ConfigDict


class MovieCreate(BaseModel):
    title: str
    description: str | None = None
    filename: str

class MovieBaseResponse(BaseModel):
    id: int
    title: str
    description: str | None
    filename: str
    created_at: datetime

    @field_serializer("created_at")
    def serialize_created_at(self, value: datetime):
        return value.isoformat()

    model_config = ConfigDict(from_attributes=True)


class MovieResponse(MovieBaseResponse):
    pass


class MovieDetailsResponse(MovieBaseResponse):
    stream_url: str


# class MovieResponse(BaseModel):
#     id: int
#     title: str
#     description: str | None
#     filename: str
#     created_at: datetime
#     # stream_url: str

#     @field_serializer("created_at")
#     def serialize_created_at(self, value: datetime):
#         return value.isoformat()

#     class Config:
#         from_attributes = True