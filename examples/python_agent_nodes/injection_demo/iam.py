"""iam — privilege management. Off-map for support-agent."""

import os

from agentfield import Agent

app = Agent(
    node_id="iam",
    agentfield_server=os.getenv("AGENTFIELD_URL", "http://localhost:8080"),
    tags=["identity-admin"],
    enable_did=True,
    vc_enabled=True,
)


@app.skill(tags=["identity-admin"])
def grant_role(agent: str, role: str) -> dict:
    return {"granted": role, "to": agent}  # privilege escalation


if __name__ == "__main__":
    app.run(auto_port=True)
