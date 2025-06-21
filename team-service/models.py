from pydantic import BaseModel
from typing import Optional

class Team(BaseModel):
    id: Optional[str] = None
    name: str
    city: str
    founded: Optional[int] = None
