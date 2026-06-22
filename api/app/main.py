from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.routers.movie_router import router as movie_router
from app.routers.auth_router import router as auth_router
from app.routers.watch_event_router import router as watch_event_router
from app.routers.watch_progress_router import router as watch_progress_router

# Will be replaced with Alembic migrations in the future, but for now we can create the tables here
# from app.database import engine
# from app.models.movie import Base

# Base.metadata.create_all(bind=engine)


app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


app.include_router(watch_progress_router)
app.include_router(watch_event_router)
app.include_router(movie_router)
app.include_router(auth_router)


@app.get("/health")
def health():
    return {
        "service": "api",
        "status": "always healthy",
    }
