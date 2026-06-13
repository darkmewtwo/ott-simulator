from pydantic import BaseModel


class MovieCreate(BaseModel):
    title: str
    description: str | None = None
    filename: str


class MovieResponse(BaseModel):
    id: int
    title: str
    description: str | None
    filename: str
    created_at: str

    class Config:
        from_attributes = True