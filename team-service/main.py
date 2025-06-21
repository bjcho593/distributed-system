from fastapi import FastAPI, HTTPException
from bson import ObjectId
from bson.errors import InvalidId

from models import Team
from database import team_collection

app = FastAPI()

@app.get("/teams/status")
def status():
    return {"message": "Team Service is running ✅"}

@app.post("/teams")
def create_team(team: Team):
    result = team_collection.insert_one(team.dict(exclude={"id"}))
    return {"id": str(result.inserted_id)}

@app.get("/teams/{team_id}")
def get_team(team_id: str):
    try:
        team = team_collection.find_one({"_id": ObjectId(team_id)})
        if not team:
            raise HTTPException(status_code=404, detail="Team not found")
        team["id"] = str(team["_id"])
        del team["_id"]
        return team
    except InvalidId:
        raise HTTPException(status_code=400, detail="Invalid team ID")
