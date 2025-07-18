from contextlib import asynccontextmanager
from fastapi import FastAPI
from app.routes.station import router as station_router
from prisma import Prisma
import time

db = Prisma()


@asynccontextmanager
async def lifespan(app: FastAPI):
    await db.connect()
    app.state.db = db
    yield
    await db.disconnect()
   
    
app = FastAPI(
    title="User Management API",
    description="API for managing users with Prisma and FastAPI",
    version="1.0.0",
    openapi_url="/openapi.json",
    openapi_tags=[
        {
            "name": "Users",
            "description": "Operations with users"
        }
    ],
    docs_url="/docs",
    redoc_url="/redoc",
    lifespan=lifespan,
    debug=True,

)

@app.middleware("http")
async def log_time(request, call_next):
    start = time.time()
    response = await call_next(request)
    duration = time.time() - start
    print(f" {request.method} {request.url.path} took {duration:.2f}s")
    return response
@app.get("/health")
async def health():
    try:
        await db.user.find_many(take=1)
        return {"status": "ok"}
    except Exception as e:
        return {"status": "error", "details": str(e)}
    
    
app.include_router(station_router, prefix="/users", tags=["Users"])
