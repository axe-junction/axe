import uuid
from pydantic import BaseModel, EmailStr
from typing import Optional

class placeCreate(BaseModel):
   Name  :str   
   Latitude:int 
   Longitude :int

class placeInDb(placeCreate):
    id: uuid.UUID
    type: str = "hybrid_bus"
    
    
    
    

