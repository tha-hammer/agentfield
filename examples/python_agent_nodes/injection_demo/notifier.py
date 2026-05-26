"""notifier — outbound channel. The exfil route, off-map for support-agent."""

import os

from agentfield import Agent

app = Agent(
    node_id="notifier",
    agentfield_server=os.getenv("AGENTFIELD_URL", "http://localhost:8080"),
    tags=["egress"],
    enable_did=True,
    vc_enabled=True,
)


@app.skill(tags=["egress"])
def send_email(to: str, body: str) -> dict:
    return {"sent_to": to, "bytes": len(body)}  # data leaves the building


if __name__ == "__main__":
    app.run(auto_port=True)
