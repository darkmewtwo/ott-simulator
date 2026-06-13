from fastapi import FastAPI
from app.routers.movie_router import router as movie_router


# Will be replaced with Alembic migrations in the future, but for now we can create the tables here
from app.database import engine
from app.models.movie import Base

Base.metadata.create_all(bind=engine)


app = FastAPI()
app.include_router(movie_router)


@app.get("/health")
def health():
    return {
        "service": "api",
        "status": "always healthy",
    }
